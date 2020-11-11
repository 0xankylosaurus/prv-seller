package agents

import (
	"fmt"
	"portalfeeders/utils"

	"github.com/0xkraken/incognito-sdk-golang/rpcclient"
	"github.com/0xkraken/incognito-sdk-golang/transaction"
	"github.com/blockcypher/gobcy"
)

func getBlockCypherAPI(networkName string) gobcy.API {
	//explicitly
	bc := gobcy.API{}
	bc.Token = "029727206f7e4c8fb19301e4629c5793"
	bc.Coin = "btc"        //options: "btc","bcy","ltc","doge"
	bc.Chain = networkName //depending on coin: "main","test3","test"
	return bc
}

func CreateAndSendNormalTx(
	client *utils.HttpClient,
	privateKeyStr string,
	paymentInfoParams map[string]uint64,
	isPrivacy bool,
) (string, error) {
	rpcClient := rpcclient.NewHttpClient(client.GetURL(), "", "", 0)
	txID, err := transaction.CreateAndSendNormalTx(rpcClient, privateKeyStr, paymentInfoParams, DefaultFee, isPrivacy)
	if err != nil {
		return "", fmt.Errorf("Error when create and send normal tx %v\n", err)
	}

	return txID, nil
}

func CreateAndSendTxRelayBNBHeader(
	client *utils.HttpClient,
	privateKeyStr string,
	bnbHeaderStr string,
	bnbBlockHeight int64,
) (string, error) {
	rpcClient := rpcclient.NewHttpClient(client.GetURL(), "", "", 0)

	txID, err := transaction.CreateAndSendTxRelayBNBHeader(rpcClient, privateKeyStr, bnbHeaderStr, bnbBlockHeight, DefaultFee)
	if err != nil {
		return "", fmt.Errorf("Error when create and send tx relay bnb block %v\n", err)
	}

	return txID, nil
}

func CreateAndSendTxRelayBTCHeader(
	client *utils.HttpClient,
	privateKeyStr string,
	btcHeaderStr string,
	btcBlockHeight int64,
) (string, error) {
	rpcClient := rpcclient.NewHttpClient(client.GetURL(), "", "", 0)

	txID, err := transaction.CreateAndSendTxRelayBTCHeader(rpcClient, privateKeyStr, btcHeaderStr, btcBlockHeight, DefaultFee)
	if err != nil {
		return "", fmt.Errorf("Error when create and send tx relay btc block %v\n", err)
	}

	return txID, nil
}

func CreateAndSendTxPortalExchangeRate(
	client *utils.HttpClient,
	privateKeyStr string,
	exchangeRateParam map[string]uint64,
) (string, error) {
	rpcClient := rpcclient.NewHttpClient(client.GetURL(), "", "", 0)

	txID, err := transaction.CreateAndSendTxPortalExchangeRate(rpcClient, privateKeyStr, exchangeRateParam, DefaultFee)
	if err != nil {
		return "", fmt.Errorf("Error when create and send tx exchange rate %v\n", err)
	}

	return txID, nil
}

// number of blocks per year: 3*60*60*24*365 = 94608000
// PRV fee for BNB relayer in 10 years: 94608000 * 10 * 20 = 18921600000 ~ 20 PRV + 10 PRV (extra)
func SplitUTXOs(
	client *utils.HttpClient,
	privateKeyStr string,
	numUTXOs int,
) error {
	rpcClient := rpcclient.NewHttpClient(client.GetURL(), "", "", 0)

	err := transaction.SplitUTXOs(rpcClient, privateKeyStr, numUTXOs, DefaultFee)
	if err != nil {
		return err
	}

	return nil
}
