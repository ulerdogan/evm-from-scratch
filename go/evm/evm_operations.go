package evm

import (
	"encoding/hex"
	"evm-from-scratch-go/utils"
	"fmt"
	"math"
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

func (m *EvmMemory) expand(offset, size int) {
	expansion := (float64(offset) + float64(size)) / 32
	expansion = math.Ceil(expansion)
	n := int(expansion) * 32

	if n > cap(m.Data) {
		d := make([]byte, n)
		copy(d[0:len(m.Data)], m.Data[:])
		m.Data = d
	}
}

func (m *EvmMemory) store(offset, size int, value *big.Int) {
	m.expand(offset, size)

	hx := utils.ToHex(value)
	if len(hx) > 64 {
		hx = hx[len(hx)-size*2:]
	} else {
		hx = fmt.Sprintf("%0*s", size*2, hx)
	}

	arr, _ := hex.DecodeString(hx)
	for i, v := range arr {
		m.Data[offset+i] = v
	}
}

func (m *EvmMemory) load(offset int) *big.Int {
	m.expand(offset, 32)

	item := ""
	for i := offset; i < offset+32; i++ {
		hx := utils.ToHex(m.Data[i])
		if len(hx) == 1 {
			hx = "0" + hx
		}
		item = item + hx
	}

	value := big.NewInt(0)
	value.SetString(item, 16)
	return value
}
