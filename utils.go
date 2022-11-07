package main

import (
	"math/big"
)

func toStrings(stack []*big.Int) []string {
	var strings []string
	for _, s := range stack {
		strings = append(strings, s.String())
	}
	return strings
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