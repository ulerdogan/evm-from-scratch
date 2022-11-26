package utils

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/big"
)

func ToStrings(stack []*big.Int) []string {
	var strings []string
	for _, s := range stack {
		strings = append(strings, s.String())
	}
	return strings
}

func ToHex(bn interface{}) string {
	return fmt.Sprintf("%x", bn)
}

func ByteToBn(b string) *big.Int {
	hx, _ := hex.DecodeString(b)
	item := ToHex(hx)
	bn := new(big.Int)
	bn.SetString(item, 16)
	return bn
}

func TwosComps(heads []*big.Int) []*big.Int {
	for i := range heads {
		heads[i] = TwosComp(heads[i])
	}
	return heads
}

func TwosComp(head *big.Int) *big.Int {
	s := fmt.Sprintf("%0*b", 256, head)
	if string(s[0]) == "1" {
		bn := flipAdd(s)
		bn.Mul(bn, big.NewInt(-1))
		head = bn
	}
	return head
}

func ConvNumber(bn *big.Int) *big.Int {
	if bn.Cmp(big.NewInt(0)) == -1 {
		s := fmt.Sprintf("%0*b", 256, bn.Mul(bn, big.NewInt(-1)))
		bn = flipAdd(s)
	}
	return bn
}

func Keccak256(bn *big.Int) *big.Int {
	h := sha3.NewLegacyKeccak256()
	h.Write(bn.Bytes())
	sum := h.Sum(nil)
	return big.NewInt(0).SetBytes(sum)
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
