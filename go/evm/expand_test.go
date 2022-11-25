package evm

import (
	"fmt"
	"testing"
)

func TestExpand(t *testing.T) {
	m := &EvmMemory{}
	m.Data = make([]byte, 32)
	m.Data[0], m.Data[1] = 1, 2
	fmt.Println(m, len(m.Data), cap(m.Data))
	m.expand(1)
	fmt.Println(m, len(m.Data), cap(m.Data))
}
