package parsers

import (
	"github.com/0xjeffro/tx-parser/solana/globals"
	"github.com/0xjeffro/tx-parser/solana/programs/pumpfun"
	"github.com/0xjeffro/tx-parser/solana/types"
	"github.com/0xjeffro/tx-parser/utils"
	"github.com/gagliardetto/solana-go"
)

type SwapData struct {
	Mint                 solana.PublicKey
	SolAmount            uint64
	TokenAmount          uint64
	IsBuy                bool
	User                 solana.PublicKey
	Timestamp            int64
	VirtualSolReserves   uint64
	VirtualTokenReserves uint64
}

func SellParser(result *types.ParsedResult, instruction types.Instruction, decodedData []byte) (action *types.PumpFunSellAction, err error) {
	datas, err := utils.CommonParseSwap[SwapData](result, instruction, pumpfun.Program)
	if err != nil {
		return nil, err
	}

	data := datas[0]
	action = &types.PumpFunSellAction{
		BaseAction: types.BaseAction{
			ProgramID:       pumpfun.Program,
			ProgramName:     pumpfun.ProgramName,
			InstructionName: "Sell",
		},
		Who:             data.User.String(),
		FromToken:       data.Mint.String(),
		ToToken:         globals.WSOL,
		FromTokenAmount: data.TokenAmount,
		ToTokenAmount:   data.SolAmount,
	}

	fromTokenDecimals, err := utils.CommonParseDecimals(result, action.FromToken)
	if err != nil {
		return nil, err
	}
	action.FromTokenDecimals = fromTokenDecimals

	toTokenDecimals, err := utils.CommonParseDecimals(result, action.ToToken)
	if err != nil {
		return nil, err
	}
	action.ToTokenDecimals = toTokenDecimals

	return action, err
}
