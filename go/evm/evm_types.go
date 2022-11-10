package evm

import "math/big"

const size uint = 2048*2048

type EvmStack struct {
	Stack []*big.Int
}

type EvmMemory struct {
	Data [size]byte
}
