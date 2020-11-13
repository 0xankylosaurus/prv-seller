package agents

const (
	// DefaultFee - default fee
	DefaultFee = 20

	// DefaultExecuteDuration - default execute duration
	DefaultExecuteDuration = 100

	PRVID   = "0000000000000000000000000000000000000000000000000000000000000004"
	PUSDTID = "716fd1009e2a1669caacc36891e707bfdf02590f96ebd897548e8963c95ebac0"
	// PUSDTID = "00000000000000000000000000000000000000000000000000000000000000ff"

	MaxSellPRVTime         = 10
	PRVRateLowerBound      = 1000000           // 1 usdt
	PRVAmountToSellAtATime = uint64(250 * 1e9) // 250 prv
	MinAcceptableAmount    = uint64(250 * 1e6) // 250 pusdt
	BurningAddress         = "12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA"
)
