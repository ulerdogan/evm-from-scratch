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

func ToAddress(bn *big.Int) string {
	return ("0x" + ToHex(bn))
}

func Contains(opcodes []byte, opcode byte) bool {
    for _, n := range opcodes {
        if n == opcode {
            return true
        }
    }
    return false
}

func HexToBn(str string) *big.Int {
	bn := new(big.Int)
	if str == "" || str == "0x" {
		return big.NewInt(0)
	}

	if str[:2] == "0x" {
		bn, _ = bn.SetString(str[2:], 16)
		return bn
	}
	bn, _ = bn.SetString(str, 16)
	return bn
}

func ByteToBn(b string) *big.Int {
	hx, _ := hex.DecodeString(b)
	item := ToHex(hx)
	bn := HexToBn(item)
	return bn
}

func PadRight(data string, amount int) string {
	for len(data) < amount {
		data = data + "0"
	}
	return data
}

func PadLeft(data string, amount int) string {
	return fmt.Sprintf("%0*s", amount, data)
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
