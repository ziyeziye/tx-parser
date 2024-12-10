package parsers

import (
	"fmt"

	"github.com/0xjeffro/tx-parser/solana/programs/orca"
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
	case orca.SwapDiscriminator:
		return SwapParser(result, instruction, decode)
	default:
		fmt.Println("Unknown discriminator", discriminator)
		return types.UnknownAction{
			BaseAction: types.BaseAction{
				ProgramID:       result.AccountList[instruction.ProgramIDIndex],
				ProgramName:     orca.ProgramName,
				InstructionName: "Unknown",
			},
		}, nil
	}
}
