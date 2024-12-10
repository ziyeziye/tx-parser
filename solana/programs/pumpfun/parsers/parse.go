package parsers

// func processPumpfunSwaps(instructionIndex int) []SwapData {
// 	var swaps []SwapData
// 	for _, innerInstructionSet := range p.tx.Meta.InnerInstructions {
// 		if innerInstructionSet.Index == uint16(instructionIndex) {
// 			for _, innerInstruction := range innerInstructionSet.Instructions {
// 				if p.isPumpFunTradeEventInstruction(innerInstruction) {
// 					eventData, err := p.parsePumpfunTradeEventInstruction(innerInstruction)
// 					if err != nil {
// 						p.Log.Errorf("error processing Pumpfun trade event: %s", err)
// 					}
// 					if eventData != nil {
// 						swaps = append(swaps, SwapData{Type: PUMP_FUN, Data: eventData})
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return swaps
// }

// func parsePumpfunTradeEventInstruction(instruction types.Instruction) (*PumpfunTradeEvent, error) {

// 	decodedBytes, err := base58.Decode(instruction.Data.String())
// 	if err != nil {
// 		return nil, fmt.Errorf("error decoding instruction data: %s", err)
// 	}

// 	var sellData SellData
// 	err := borsh.Deserialize(&sellData, decodedData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	decoder := borsh.Deserialize(decodedBytes[16:])

// 	return handlePumpfunTradeEvent(decoder)
// }
