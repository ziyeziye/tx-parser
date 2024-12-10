package types

import "github.com/gagliardetto/solana-go"

type ParsedResult struct {
	RawTx           RawTx    `json:"rawTx"`
	AccountList     []string `json:"accountList"`
	Actions         []Action `json:"actions"`
	SplTokenInfoMap map[string]TokenInfo
	SplDecimalsMap  map[string]uint8
}

type TokenInfo struct {
	Mint     string
	Decimals uint8
}

type Action interface {
	GetProgramID() string
	GetProgramName() string
	GetInstructionName() string
}

type BaseAction struct {
	ProgramID       string `json:"programId"`
	ProgramName     string `json:"programName"`
	InstructionName string `json:"instructionName"`
}

func (a BaseAction) GetProgramID() string {
	return a.ProgramID
}

func (a BaseAction) GetProgramName() string {
	return a.ProgramName
}

func (a BaseAction) GetInstructionName() string {
	return a.InstructionName
}

type UnknownAction struct {
	BaseAction
	Error error `json:"error"`
}

func (p *ParsedResult) ExtractSPLTokenInfo() error {
	splTokenAddresses := make(map[string]TokenInfo)

	for _, accountInfo := range p.RawTx.Meta.PostTokenBalances {
		if !GetPublicKey(accountInfo.Mint).IsZero() {
			accountKey := p.AccountList[accountInfo.AccountIndex]
			splTokenAddresses[accountKey] = TokenInfo{
				Mint:     accountInfo.Mint,
				Decimals: accountInfo.UITokenAmount.Decimals,
			}
		}
	}

	processInstruction := func(instr solana.CompiledInstruction) {
		if !GetPublicKey(p.AccountList[instr.ProgramIDIndex]).Equals(solana.TokenProgramID) {
			return
		}

		if len(instr.Data) == 0 || (instr.Data[0] != 3 && instr.Data[0] != 12) {
			return
		}

		if len(instr.Accounts) < 3 {
			return
		}

		source := p.AccountList[instr.Accounts[0]]
		destination := p.AccountList[instr.Accounts[1]]

		if _, exists := splTokenAddresses[source]; !exists {
			splTokenAddresses[source] = TokenInfo{Mint: "", Decimals: 0}
		}
		if _, exists := splTokenAddresses[destination]; !exists {
			splTokenAddresses[destination] = TokenInfo{Mint: "", Decimals: 0}
		}
	}

	for _, instr := range p.RawTx.Transaction.Message.Instructions {
		processInstruction(instr)
	}
	for _, innerSet := range p.RawTx.Meta.InnerInstructions {
		for _, instr := range innerSet.Instructions {
			processInstruction(instr)
		}
	}

	for account, info := range splTokenAddresses {
		if info.Mint == "" {
			splTokenAddresses[account] = TokenInfo{
				Mint:     solana.SolMint.String(),
				Decimals: 9, // Native SOL has 9 decimal places
			}
		}
	}

	p.SplTokenInfoMap = splTokenAddresses

	return nil
}

func (p *ParsedResult) ExtractSPLDecimals() error {
	mintToDecimals := make(map[string]uint8)
	// var tokenBalances []TokenBalance
	// tokenBalances = append(tokenBalances, p.RawTx.Meta.PreTokenBalances...)
	// tokenBalances = append(tokenBalances, p.RawTx.Meta.PostTokenBalances...)
	// for _, v := range tokenBalances {
	// 	mintToDecimals[v.Mint] = v.UITokenAmount.Decimals
	// }

	for _, accountInfo := range p.RawTx.Meta.PostTokenBalances {
		if !GetPublicKey(accountInfo.Mint).IsZero() {
			mintAddress := accountInfo.Mint
			mintToDecimals[mintAddress] = uint8(accountInfo.UITokenAmount.Decimals)
		}
	}

	processInstruction := func(instr solana.CompiledInstruction) {
		if !GetPublicKey(p.AccountList[instr.ProgramIDIndex]).Equals(solana.TokenProgramID) {
			return
		}

		if len(instr.Data) == 0 || (instr.Data[0] != 3 && instr.Data[0] != 12) {
			return
		}

		if len(instr.Accounts) < 3 {
			return
		}

		mint := p.AccountList[instr.Accounts[1]]
		if _, exists := mintToDecimals[mint]; !exists {
			mintToDecimals[mint] = 0
		}
	}

	for _, instr := range p.RawTx.Transaction.Message.Instructions {
		processInstruction(instr)
	}
	for _, innerSet := range p.RawTx.Meta.InnerInstructions {
		for _, instr := range innerSet.Instructions {
			processInstruction(instr)
		}
	}

	// Add Native SOL if not present
	if _, exists := mintToDecimals[solana.SolMint.String()]; !exists {
		mintToDecimals[solana.SolMint.String()] = 9 // Native SOL has 9 decimal places
	}

	p.SplDecimalsMap = mintToDecimals

	return nil
}
