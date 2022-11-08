package evm

import "math/big"

type EvmStack struct {
	S []*big.Int
}