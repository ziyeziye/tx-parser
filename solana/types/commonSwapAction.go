package types

type CommonSwapAction struct {
	BaseAction
	Who               string `json:"who"`
	FromToken         string `json:"fromToken"`
	FromTokenAmount   uint64 `json:"fromTokenAmount"`
	FromTokenDecimals uint8  `json:"fromTokenDecimals"`
	ToToken           string `json:"toToken"`
	ToTokenAmount     uint64 `json:"toTokenAmount"`
	ToTokenDecimals   uint8  `json:"toTokenDecimals"`
}
