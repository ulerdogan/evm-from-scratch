package evm

import "math/big"

type EvmStack struct {
	Stack []*big.Int
}

type EvmMemory struct {
	Data []byte
}
