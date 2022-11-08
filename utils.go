package main

import (
	"fmt"
	"math/big"
)

func toStrings(stack []*big.Int) []string {
	var strings []string
	for _, s := range stack {
		strings = append(strings, s.String())
	}
	return strings
}

func twosComps(heads []*big.Int) []*big.Int {
	heads[0] = twosComp(heads[0])
	heads[1] = twosComp(heads[1])
	return heads
}

func twosComp(head *big.Int) *big.Int {
	s := fmt.Sprintf("%0*b", 256, head)
	if string(s[0]) == "1" {
		bn := flipAdd(s)
		bn.Mul(bn, big.NewInt(-1))
		head = bn
	}
	return head
}

func convNumber(bn *big.Int) *big.Int {
	if bn.Cmp(big.NewInt(0)) == -1 {
		s := fmt.Sprintf("%0*b", 256, bn.Mul(bn, big.NewInt(-1)))
		bn = flipAdd(s)
	}
	return bn
}

func flipAdd(s string) *big.Int {
	str := flipStr(s)
	bn, _ := new(big.Int).SetString(str, 2)
	bn.Add(bn, big.NewInt(1))
	return bn
}

func flipStr(b string) string {
	str := ""
	for _, c := range b {
		switch string(c) {
		case "0":
			str = str + "1"
		case "1":
			str = str + "0"
		}
	}
	return str
}
