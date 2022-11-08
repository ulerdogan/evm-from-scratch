package main

import (
	"math/big"
)

func (s *evmstack) getHeads() []*big.Int {
	heads := []*big.Int{s.s[0], s.s[1]}
	s.s = s.s[2:]
	return heads
}

func (s *evmstack) getHead() *big.Int {
	head := s.s[0]
	s.s = s.s[1:]
	return head
}

func (s *evmstack) oprHeads(f func(x *big.Int, y *big.Int) *big.Int, signed bool) *big.Int {
	heads := s.getHeads()
	if signed {
		heads = twosComps(heads)
		res := f(heads[0], heads[1])
		return convNumber(res)
	}
	return f(heads[0], heads[1])
}

func (s *evmstack) oprHead(f func(x *big.Int) *big.Int, signed bool) *big.Int {
	head := s.getHead()
	if signed {
		head = twosComp(head)
		res := f(head)
		return convNumber(res)
	}
	return f(head)
}
