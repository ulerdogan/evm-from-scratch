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
	for i := big.NewInt(0); i.Cmp(big.NewInt(32)) == -1; i.Add(i, big.NewInt(1)) {
		d := int(i.Add(i, offset).Int64())
		v := new(big.Int).Sub(big.NewInt(32), i)
		v.Sub(v, big.NewInt(1)).Mul(v, big.NewInt(8))
		m.Data[d] = uint8(value.Rsh(value, uint(v.Int64())).And(value, utils.ByteToBn("ff")).Int64())
	}
}

func (m *EvmMemory) load(offset *big.Int) *big.Int {
	value := big.NewInt(0)
	for i := big.NewInt(0); i.Cmp(big.NewInt(32)) == -1; i.Add(i, big.NewInt(1)) {
		o := big.NewInt(int64(m.Data[int((offset.Add(offset, i)).Int64())]))
		value.Lsh(value, 8).Or(value, o)
	}
	return value
}
