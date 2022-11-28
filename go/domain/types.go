package domain

import (
	"math/big"
)

type TestCase struct {
	Name   string
	State  map[string]AccState
	Block  BlockInfo
	Tx     TxData
	Code   code
	Expect expect
}

type code struct {
	Bin string
	Asm string
}

type AccState struct {
	Balance string
	Nonce   int
	Code    code
}

type BlockInfo struct {
	Coinbase   string
	Timestamp  string
	Number     string
	Difficulty string
	GasLimit   string
	ChainId    string
}

type TxData struct {
	From     string
	To       string
	Origin   string
	GasPrice string
	Value    string
	Data     string
}

type Logs struct {
	Address string
	Data string
	Topics []string
}

type Result struct {
	Stack []*big.Int
	Return string
	Success bool
}

type maxNums struct {
	Uint256Max *big.Int
}

type expect struct {
	Stack   []string
	Success bool
	Return  string
}

var (
	Max *maxNums = &maxNums{
		Uint256Max: func() *big.Int {
			var Max big.Int
			Max.Exp(big.NewInt(2), big.NewInt(256), nil)
			return &Max
		}(),
	}
)
