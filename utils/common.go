package utils

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/0xjeffro/tx-parser/solana/globals"
	"github.com/0xjeffro/tx-parser/solana/types"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

func CommonParseSwap[T any](result *types.ParsedResult, instruction types.Instruction, programID string, innerProgramIds ...string) ([]*T, error) {
	var instructionIndex int
	for idx, instr := range result.RawTx.Transaction.Message.Instructions {
		if result.AccountList[instr.ProgramIDIndex] == programID && instr.Data.String() == instruction.Data.String() {
			instructionIndex = idx
			break
		}
	}

	var instructions []types.Instruction
	for _, innerInstruction := range result.RawTx.Meta.InnerInstructions {
		if innerInstruction.Index == instructionIndex {
			instructions = innerInstruction.Instructions
			break
		}
	}

	innerProgramId := programID
	if len(innerProgramIds) > 0 {
		innerProgramId = innerProgramIds[0]
	}

	actions := make([]*T, 0)
	for _, instr := range instructions {
		programId := result.AccountList[instr.ProgramIDIndex]
		switch programId {
		case innerProgramId:
			action := new(T)
			decodedBytes, err := base58.Decode(instr.Data.String())
			if err != nil {
				log.Printf("error decoding instruction data: %s", err)
				continue
			}
			if len(decodedBytes) < 16 {
				continue
			}
			decoder := bin.NewBorshDecoder(decodedBytes[16:])
			// var data SwapData
			err = decoder.Decode(action)
			if err != nil {
				log.Printf("error decoding instruction data: %s", err)
				continue
			}
			actions = append(actions, action)
		default:
			fmt.Printf("commonParse: Program:%s, InnerProgram:%s unknown inner program:%s", programID, innerProgramId, programId)
		}
	}

	if len(actions) > 0 {
		return actions, nil
	}

	return nil, fmt.Errorf("unknown instruction")
}

func CommonParseTransfer(result *types.ParsedResult, instruction types.Instruction, programID string) ([]*TransferData, error) {
	var instructionIndex int
	for idx, instr := range result.RawTx.Transaction.Message.Instructions {
		if result.AccountList[instr.ProgramIDIndex] == programID && instr.Data.String() == instruction.Data.String() {
			instructionIndex = idx
			break
		}
	}

	var instructions []types.Instruction
	for _, innerInstruction := range result.RawTx.Meta.InnerInstructions {
		if innerInstruction.Index == instructionIndex {
			instructions = innerInstruction.Instructions
			break
		}
	}

	actions := make([]*TransferData, 0)
	for _, instr := range instructions {
		programId := result.AccountList[instr.ProgramIDIndex]
		switch programId {
		case solana.TokenProgramID.String():
			action := new(TransferData)
			decodedBytes, err := base58.Decode(instr.Data.String())
			if err != nil {
				log.Printf("error decoding instruction data: %s", err)
				continue
			}

			switch {
			case isTransfer(result.AccountList, solana.TokenProgramID, instr.Accounts, decodedBytes):
				action = processTransfer(result, instr, decodedBytes)
			case isTransferCheck(result.AccountList, solana.TokenProgramID, instr.Accounts, decodedBytes):
				check := processTransferCheck(result, instr, decodedBytes)
				amount, _ := strconv.Atoi(check.Info.TokenAmount.Amount)
				action = &TransferData{
					Info: TransferInfo{
						Amount:      uint64(amount),
						Source:      check.Info.Source,
						Destination: check.Info.Destination,
						Authority:   check.Info.Authority,
					},
					Type:     check.Type,
					Mint:     check.Info.Mint,
					Decimals: check.Info.TokenAmount.Decimals,
				}
			default:
				fmt.Printf("commonParseTransfer: Program:%s, unknown inner program:%s", programID, programId)
			}

			actions = append(actions, action)
		default:
			fmt.Printf("commonParseTransfer: Program:%s, unknown inner program:%s", programID, programId)
		}
	}

	if len(actions) > 0 {
		return actions, nil
	}

	return nil, fmt.Errorf("unknown instruction")
}

type TransferInfo struct {
	Amount      uint64 `json:"amount"`
	Authority   string `json:"authority"`
	Destination string `json:"destination"`
	Source      string `json:"source"`
}

type TransferData struct {
	Info     TransferInfo `json:"info"`
	Type     string       `json:"type"`
	Mint     string       `json:"mint"`
	Decimals uint8        `json:"decimals"`
}

type TransferCheck struct {
	Info struct {
		Authority   string `json:"authority"`
		Destination string `json:"destination"`
		Mint        string `json:"mint"`
		Source      string `json:"source"`
		TokenAmount struct {
			Amount         string  `json:"amount"`
			Decimals       uint8   `json:"decimals"`
			UIAmount       float64 `json:"uiAmount"`
			UIAmountString string  `json:"uiAmountString"`
		} `json:"tokenAmount"`
	} `json:"info"`
	Type string `json:"type"`
}

func processTransfer(result *types.ParsedResult, instr types.Instruction, decodedBytes []byte) *TransferData {
	amount := binary.LittleEndian.Uint64(decodedBytes[1:9])

	transferData := &TransferData{
		Info: TransferInfo{
			Amount:      amount,
			Source:      result.AccountList[instr.Accounts[0]],
			Destination: result.AccountList[instr.Accounts[1]],
			Authority:   result.AccountList[instr.Accounts[2]],
		},
		Type:     "transfer",
		Mint:     result.SplTokenInfoMap[result.AccountList[instr.Accounts[1]]].Mint,
		Decimals: result.SplTokenInfoMap[result.AccountList[instr.Accounts[1]]].Decimals,
	}

	if transferData.Mint == "" {
		transferData.Mint = "Unknown"
	}

	return transferData
}

func processTransferCheck(result *types.ParsedResult, instr types.Instruction, decodedBytes []byte) *TransferCheck {

	amount := binary.LittleEndian.Uint64(decodedBytes[1:9])

	transferData := &TransferCheck{
		Type: "transferChecked",
	}

	transferData.Info.Source = result.AccountList[instr.Accounts[0]]
	transferData.Info.Destination = result.AccountList[instr.Accounts[2]]
	transferData.Info.Mint = result.AccountList[instr.Accounts[1]]
	transferData.Info.Authority = result.AccountList[instr.Accounts[3]]

	transferData.Info.TokenAmount.Amount = fmt.Sprintf("%d", amount)
	transferData.Info.TokenAmount.Decimals = result.SplTokenInfoMap[transferData.Info.Mint].Decimals
	uiAmount := float64(amount) / math.Pow10(int(transferData.Info.TokenAmount.Decimals))
	transferData.Info.TokenAmount.UIAmount = uiAmount
	transferData.Info.TokenAmount.UIAmountString = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.9f", uiAmount), "0"), ".")

	return transferData
}

// isTransfer checks if the instruction is a token transfer (Raydium, Orca)
func isTransfer(allAccountKeys []string, progID solana.PublicKey, innerAccounts []uint16, decodedBytes []byte) bool {
	if !progID.Equals(solana.TokenProgramID) {
		return false
	}

	if len(innerAccounts) < 3 || len(decodedBytes) < 9 {
		return false
	}

	if decodedBytes[0] != 3 {
		return false
	}

	for i := 0; i < 3; i++ {
		if int(innerAccounts[i]) >= len(allAccountKeys) {
			return false
		}
	}

	return true
}

// isTransferCheck checks if the instruction is a token transfer check (Meteora)
func isTransferCheck(allAccountKeys []string, progID solana.PublicKey, innerAccounts []uint16, decodedBytes []byte) bool {
	if !progID.Equals(solana.TokenProgramID) && !progID.Equals(solana.Token2022ProgramID) {
		return false
	}

	if len(innerAccounts) < 4 || len(decodedBytes) < 9 {
		return false
	}

	if decodedBytes[0] != 12 {
		return false
	}

	for i := 0; i < 4; i++ {
		if int(innerAccounts[i]) >= len(allAccountKeys) {
			return false
		}
	}

	return true
}

func CommonParseDecimals(result *types.ParsedResult, mint string) (uint8, error) {
	if mint == globals.WSOL {
		return globals.SOLDecimals, nil
	}

	decimals := result.SplDecimalsMap[mint]
	if decimals > 0 {
		return decimals, nil
	}

	return 0, fmt.Errorf("unknown mint")
}
