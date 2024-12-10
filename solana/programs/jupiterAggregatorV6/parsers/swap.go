package parsers

import (
	"github.com/0xjeffro/tx-parser/solana/programs/jupiterAggregatorV6"
	"github.com/0xjeffro/tx-parser/solana/types"
	"github.com/0xjeffro/tx-parser/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/samber/lo"
)

type SwapData struct {
	Amm          solana.PublicKey `json:"amm"`
	InputMint    solana.PublicKey `json:"inputMint"`
	InputAmount  uint64           `json:"inputAmount"`
	OutputMint   solana.PublicKey `json:"outputMint"`
	OutputAmount uint64           `json:"outputAmount"`
}

func SwapParser(result *types.ParsedResult, instruction types.Instruction) (*types.JupiterAggregatorV6RouteAction, error) {

	datas, err := utils.CommonParseSwap[SwapData](result, instruction, jupiterAggregatorV6.Program)
	if err != nil {
		return nil, err
	}

	data := datas[0]
	last := lo.LastOrEmpty(datas)
	data.OutputMint = last.OutputMint
	data.OutputAmount = last.OutputAmount

	action := &types.JupiterAggregatorV6RouteAction{
		BaseAction: types.BaseAction{
			ProgramID:       jupiterAggregatorV6.Program,
			ProgramName:     jupiterAggregatorV6.ProgramName,
			InstructionName: "Swap",
		},
		Who:             result.AccountList[instruction.Accounts[0]],
		FromToken:       data.InputMint.String(),
		ToToken:         data.OutputMint.String(),
		FromTokenAmount: data.InputAmount,
		ToTokenAmount:   data.OutputAmount,
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
