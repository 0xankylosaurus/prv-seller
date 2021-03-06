package agents

import (
	"errors"
	"fmt"
	"github.com/0xkraken/incognito-sdk-golang/rpcclient"
	"github.com/0xkraken/incognito-sdk-golang/transaction"
	"math/big"
	"portalfeeders/entities"
	"time"
)

type PRVSeller struct {
	AgentAbs
	Counter       uint
	SellerPrivKey string
	SellerAddress string
}

func (b *PRVSeller) getLatestBeaconHeight() (uint64, error) {
	params := []interface{}{}
	var beaconBestStateRes entities.BeaconBestStateRes
	err := b.RPCClient.RPCCall("getbeaconbeststate", params, &beaconBestStateRes)
	if err != nil {
		return 0, err
	}

	if beaconBestStateRes.RPCError != nil {
		b.Logger.Errorf("getLatestBeaconHeight: call RPC error, %v\n", beaconBestStateRes.RPCError.StackTrace)
		return 0, errors.New(beaconBestStateRes.RPCError.Message)
	}
	return beaconBestStateRes.Result.BeaconHeight, nil
}

func (b *PRVSeller) getPDEState(beaconHeight uint64) (*entities.PDEState, error) {
	params := []interface{}{
		map[string]uint64{
			"BeaconHeight": beaconHeight,
		},
	}
	var pdeStateRes entities.PDEStateRes
	err := b.RPCClient.RPCCall("getpdestate", params, &pdeStateRes)
	if err != nil {
		return nil, err
	}

	if pdeStateRes.RPCError != nil {
		b.Logger.Errorf("getPDEState: call RPC error, %v\n", pdeStateRes.RPCError.StackTrace)
		return nil, errors.New(pdeStateRes.RPCError.Message)
	}
	return pdeStateRes.Result, nil
}

func (b *PRVSeller) getPRVRate() (uint64, error) {
	latestBeaconHeight, err := b.getLatestBeaconHeight()
	if err != nil {
		return 0, err
	}
	pdeState, err := b.getPDEState(latestBeaconHeight)
	if err != nil {
		return 0, err
	}
	poolPairs := pdeState.PDEPoolPairs
	prvPustPairKey := fmt.Sprintf("pdepool-%d-%s-%s", latestBeaconHeight, PRVID, PUSDTID)
	prvPustPair, found := poolPairs[prvPustPairKey]
	if !found || prvPustPair.Token1PoolValue == 0 || prvPustPair.Token2PoolValue == 0 {
		return 0, nil
	}

	tokenPoolValueToBuy := prvPustPair.Token1PoolValue
	tokenPoolValueToSell := prvPustPair.Token2PoolValue
	if prvPustPair.Token1IDStr == PRVID {
		tokenPoolValueToSell = prvPustPair.Token1PoolValue
		tokenPoolValueToBuy = prvPustPair.Token2PoolValue
	}

	invariant := big.NewInt(0)
	invariant.Mul(big.NewInt(int64(tokenPoolValueToSell)), big.NewInt(int64(tokenPoolValueToBuy)))
	newTokenPoolValueToSell := big.NewInt(0)
	newTokenPoolValueToSell.Add(big.NewInt(int64(tokenPoolValueToSell)), big.NewInt(int64(1e9)))

	newTokenPoolValueToBuy := big.NewInt(0).Div(invariant, newTokenPoolValueToSell).Uint64()
	modValue := big.NewInt(0).Mod(invariant, newTokenPoolValueToSell)
	if modValue.Cmp(big.NewInt(0)) != 0 {
		newTokenPoolValueToBuy++
	}
	if tokenPoolValueToBuy <= newTokenPoolValueToBuy {
		return 0, nil
	}
	return tokenPoolValueToBuy - newTokenPoolValueToBuy, nil
}

func (b *PRVSeller) sellPRV(sellAmt uint64) (string, error) {
	//params := []interface{}{
	//	b.SellerPrivKey,
	//	map[string]string{
	//		BurningAddress: fmt.Sprintf("%d", sellAmt),
	//	},
	//	100,
	//	-1,
	//	map[string]string{
	//		"TokenIDToBuyStr":     PUSDTID,
	//		"TokenIDToSellStr":    PRVID,
	//		"SellAmount":          fmt.Sprintf("%d", sellAmt),
	//		"MinAcceptableAmount": fmt.Sprintf("%d", MinAcceptableAmount),
	//		"TradingFee":          "0",
	//		"TraderAddressStr":    b.SellerAddress,
	//	},
	//}
	//fmt.Println("huhu params: ", params)
	//var prvTradeRes entities.PRVTradeRes
	//err := b.RPCClient.RPCCall("createandsendtxwithprvcrosspooltradereq", params, &prvTradeRes)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if prvTradeRes.RPCError != nil {
	//	b.Logger.Errorf("prvTrade: call RPC error, %v\n", prvTradeRes.RPCError.StackTrace)
	//	return nil, errors.New(prvTradeRes.RPCError.Message)
	//}
	//return &prvTradeRes, nil

	rpcClient := rpcclient.NewHttpClient(b.RPCClient.GetURL(), "", "", 0)
	txID, err := transaction.CreateAndSendTxPRVCrossPoolTrade(
		rpcClient,
		b.SellerPrivKey,
		PUSDTID,
		sellAmt,
		MinAcceptableAmount,
		0, 100)

	if err != nil {
		b.Logger.Errorf("prvTrade: Create tx error, %v\n", err)
		return "", err
	}

	return txID, nil
}

func (b *PRVSeller) Execute() {
	now := time.Now()
	hr := now.Hour()
	if b.Counter == MaxSellPRVTime && hr == 0 { // reset counter at 0 oclock
		b.Counter = 0
	}

	if b.Counter == MaxSellPRVTime {
		b.Logger.Info("Reached to max number of selling prv!")
		return
	}

	b.Logger.Info("PRVSeller agent is executing...")
	prvRate, err := b.getPRVRate()
	if err != nil {
		msg := fmt.Sprintf("PRVSeller: has a error, %v\n", err)
		b.Logger.Errorf(msg)
		// utils.SendSlackNotification(msg)
		return
	}
	fmt.Println("haha prv rate: ", prvRate)

	if prvRate < PRVRateLowerBound {
		b.Logger.Infof("PRV rate is lower the expected bound: %d", prvRate)
		return
	}

	txID, err := b.sellPRV(PRVAmountToSellAtATime)
	if err != nil {
		b.Logger.Errorf("sell prv failed with error: %v", err)
		return
	}

	b.Logger.Infof("sell prv successfully with tx: %s", txID)

	b.Counter++
	b.Logger.Info("PRVSeller agent finished...")
}
