package parsers

import (
	"fmt"

	"github.com/0xjeffro/tx-parser/solana/programs/jupiterAggregatorV6"
	"github.com/0xjeffro/tx-parser/solana/types"
	"github.com/mr-tron/base58"
)

func InstructionRouter(result *types.ParsedResult, instruction types.Instruction) (types.Action, error) {
	data := instruction.Data
	decode, err := base58.Decode(data.String())
	if err != nil {
		return nil, err
	}
	discriminator := *(*[8]byte)(decode[:8])

	switch discriminator {
	case jupiterAggregatorV6.SharedAccountsRouteDiscriminator:
		return SharedAccountsRouteParser(result, instruction)
	case jupiterAggregatorV6.RouteDiscriminator:
		return RouteParser(result, instruction)
	case jupiterAggregatorV6.ExactOutRoute, jupiterAggregatorV6.SharedAccountsExactOutRoute:
		return SwapParser(result, instruction)
	default:
		fmt.Println("Jupiter Unknown discriminator", discriminator)
		return types.UnknownAction{
			BaseAction: types.BaseAction{
				ProgramID:       result.AccountList[instruction.ProgramIDIndex],
				ProgramName:     jupiterAggregatorV6.ProgramName,
				InstructionName: "Unknown",
			},
		}, nil
	}
}
