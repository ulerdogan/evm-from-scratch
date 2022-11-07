package main

import (
	"fmt"
	"math/big"
)

func (s *evmstack) getHeads() []*big.Int {
	heads := []*big.Int{s.s[0], s.s[1]}
	s.s = s.s[2:]
	return heads
}

func (s *evmstack) oprHeads(f func(x *big.Int, y *big.Int) *big.Int, signed bool) *big.Int {
	heads := []*big.Int{s.s[0], s.s[1]}
	s.s = s.s[2:]
	if signed {
		heads = twosComp(heads)
		res := f(heads[0], heads[1])
		if res.Cmp(big.NewInt(0)) == -1 {
			s := fmt.Sprintf("%0*b", 256, res.Mul(res, big.NewInt(-1)))
			bn := flipAdd(s)
			res = bn
		}
		return res
	}
	return f(heads[0], heads[1])
}
