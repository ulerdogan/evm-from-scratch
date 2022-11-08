package main

import (
	"math/big"
)

func (s *evmstack) getHeads(n int) []*big.Int {
	heads := s.s[0:n]
	s.s = s.s[n:]
	return heads
}

func (s *evmstack) oprHeads(f func(x *big.Int, y *big.Int) *big.Int, signed bool) *big.Int {
	heads := s.getHeads(2)
	if signed {
		heads = twosComps(heads)
		res := f(heads[0], heads[1])
		return convNumber(res)
	}
	return f(heads[0], heads[1])
}

func (s *evmstack) oprHead(f func(x *big.Int) *big.Int, signed bool) *big.Int {
	head := s.getHeads(1)[0]
	if signed {
		head = twosComp(head)
		res := f(head)
		return convNumber(res)
	}
	return f(head)
}
