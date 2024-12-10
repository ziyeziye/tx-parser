package parsers

import (
	"github.com/0xjeffro/tx-parser/solana/globals"
	"github.com/0xjeffro/tx-parser/solana/programs/pumpfun"
	"github.com/0xjeffro/tx-parser/solana/types"
	"github.com/0xjeffro/tx-parser/utils"
)

func BuyParser(result *types.ParsedResult, instruction types.Instruction, decodedData []byte) (*types.PumpFunBuyAction, error) {
	datas, err := utils.CommonParseSwap[SwapData](result, instruction, pumpfun.Program)
	if err != nil {
		return nil, err
	}

	data := datas[0]
	action := &types.PumpFunBuyAction{
		BaseAction: types.BaseAction{
			ProgramID:       pumpfun.Program,
			ProgramName:     pumpfun.ProgramName,
			InstructionName: "Buy",
		},
		Who:             data.User.String(),
		FromToken:       globals.WSOL,
		ToToken:         data.Mint.String(),
		FromTokenAmount: data.SolAmount,
		ToTokenAmount:   data.TokenAmount,
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

	return action, nil
}
