package utils

import (
	"fmt"
	"math/big"
	"testing"
)

func TestByteToBn(t *testing.T) {
	tst := ByteToBn("ff")
	wnt := big.NewInt(255)

	if tst.Cmp(wnt) != int(0) {
		t.Errorf("Unmatched ops. Want %v, get %v", wnt, tst)
	}
}

func TestKeccak256(t *testing.T) {
	bn := big.NewInt(123124124)
	hash := Keccak256(bn)
	fmt.Println("hash:", hash)
}
