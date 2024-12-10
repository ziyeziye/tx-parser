package parsers

import (
	"github.com/0xjeffro/tx-parser/solana/programs/orca"
	"github.com/0xjeffro/tx-parser/solana/types"
	"github.com/0xjeffro/tx-parser/utils"
	"github.com/samber/lo"
)

func SwapParser(result *types.ParsedResult, instruction types.Instruction, decodedData []byte) (*types.CommonSwapAction, error) {

	datas, err := utils.CommonParseTransfer(result, instruction, orca.Program)
	if err != nil {
		return nil, err
	}

	data := datas[0]
	last := lo.LastOrEmpty(datas)

	action := &types.CommonSwapAction{
		BaseAction: types.BaseAction{
			ProgramID:       orca.Program,
			ProgramName:     orca.ProgramName,
			InstructionName: "Swap",
		},
		Who:             result.AccountList[instruction.Accounts[2]],
		FromToken:       data.Mint,
		ToToken:         last.Mint,
		FromTokenAmount: data.Info.Amount,
		ToTokenAmount:   last.Info.Amount,
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
