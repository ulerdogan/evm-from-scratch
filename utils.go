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

func twosComp(heads []*big.Int) []*big.Int {
	s := fmt.Sprintf("%0*b", 256, heads[0])
	if string(s[0]) == "1" {
		bn := flipAdd(s)
		bn.Mul(bn, big.NewInt(-1))
		heads[0] = bn
	}

	s = fmt.Sprintf("%0*b", 256, heads[1])
	if string(s[0]) == "1" {
		bn := flipAdd(s)
		bn.Mul(bn, big.NewInt(-1))
		heads[1] = bn
	}
	return heads
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
