package evm

import (
	"evm-from-scratch-go/utils"
	"math/big"
)

func (s *EvmStack) getHeads(n int) []*big.Int {
	if !s.checkStack(n) {
		return nil
	}
	heads := s.S[0:n]
	s.S = s.S[n:]
	return heads
}

func (s *EvmStack) oprHeads(f func(x *big.Int, y *big.Int) *big.Int, signed bool) *big.Int {
	heads := s.getHeads(2)
	if heads == nil {
		return nil
	}

	if signed {
		heads = utils.TwosComps(heads)
		res := f(heads[0], heads[1])
		return utils.ConvNumber(res)
	}
	return f(heads[0], heads[1])
}

func (s *EvmStack) oprHead(f func(x *big.Int) *big.Int, signed bool) *big.Int {
	head := s.getHeads(1)[0]
	if head == nil {
		return nil
	}

	if signed {
		head = utils.TwosComp(head)
		res := f(head)
		return utils.ConvNumber(res)
	}
	return f(head)
}

func (s *EvmStack) checkStack(n int) bool {
	return len(s.S) >= n
}
