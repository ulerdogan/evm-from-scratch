package utils

import (
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