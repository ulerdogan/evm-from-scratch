package evm

import (
	"evm-from-scratch-go/utils"
	"math/big"
)

func (s *EvmStack) getHeads(n int) []*big.Int {
	if !s.checkStack(n) {
		return nil
	}
	heads := s.Stack[0:n]
	s.Stack = s.Stack[n:]
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
	return len(s.Stack) >= n
}

func (m *EvmMemory) store(offset, value *big.Int) {
	for i := 0; i < 32; i++ {
		k := i + int(offset.Int64())
		rsh := new(big.Int).Rsh(value, uint(8*(31-i)))
		v := rsh.And(rsh, utils.ByteToBn("ff")).Int64()
		m.Data[k] = byte(v)
	}
}

func (m *EvmMemory) load(offset *big.Int) *big.Int {
	value := big.NewInt(0)
	for i := 0; i < 32; i++ {
		value.Or(value.Lsh(value, 8), big.NewInt(int64(m.Data[int(offset.Int64())+i])))
	}
	return value
}
