package types

type PumpFunBuyAction = CommonSwapAction

type PumpFunSellAction = CommonSwapAction

type PumpFunCreateAction struct {
	BaseAction
	Who                    string `json:"who"`
	Mint                   string `json:"mint"`
	MintAuthority          string `json:"mintAuthority"`
	BondingCurve           string `json:"bondingCurve"`
	AssociatedBondingCurve string `json:"associatedBondingCurve"`
	MplTokenMetadata       string `json:"mplTokenMetadata"`
	MetaData               string `json:"metaData"`

	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Uri    string `json:"uri"`
}
