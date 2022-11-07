package main

import (
	"math/big"
)

type TestCase struct {
	Name   string
	Code   code
	Expect expect
}

type code struct {
	Bin string
	Asm string
}

type evmstack struct {
	s []*big.Int
}

type maxNums struct {
	uint256Max *big.Int
}

type expect struct {
	Stack   []string
	Success bool
	Return  string
}

var (
	max *maxNums = &maxNums{
		uint256Max: func() *big.Int {
			var max big.Int
			max.Exp(big.NewInt(2), big.NewInt(256), nil)
			return &max
		}(),
	}
)
