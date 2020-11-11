package entities

type PDEPoolPair struct {
	Token1IDStr     string
	Token1PoolValue uint64
	Token2IDStr     string
	Token2PoolValue uint64
}

type PDEState struct {
	PDEPoolPairs map[string]*PDEPoolPair
}

type PDEStateRes struct {
	RPCBaseRes
	Result *PDEState
}

type PRVTrade struct {
	TxID string
}

type PRVTradeRes struct {
	RPCBaseRes
	Result *PRVTrade
}
