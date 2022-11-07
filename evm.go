/**
 * EVM From Scratch
 * Go template
 *
 * To work on EVM From Scratch in Go:
 *
 * - Install Golang: https://golang.org/doc/install
 * - Go to the `go` directory: `cd go`
 * - Edit `evm.go` (this file!), see TODO below
 * - Run `go run evm.go` to run the tests
 */

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
)

type code struct {
	Bin string
	Asm string
}

type evmstack struct {
	s []*big.Int
}

type expect struct {
	Stack   []string
	Success bool
	Return  string
}
type maxNums struct {
	uint256Max *big.Int
}

type TestCase struct {
	Name   string
	Code   code
	Expect expect
}

var (
	max *maxNums = &maxNums{
		uint256Max: func() *big.Int {
			var max big.Int
			max.Exp(big.NewInt(2), big.NewInt(256), nil)
			return &max
		}(),
	}
)

func (s *evmstack) oprHeads(f func(x *big.Int, y *big.Int) *big.Int, signed bool) *big.Int {
	heads := []*big.Int{s.s[0], s.s[1]}
	s.s = s.s[2:]
	if signed {
		s := fmt.Sprintf("%0*b", 256, heads[0])
		if string(s[0]) == "1" {
			str := flipStr(s)
			bn, _ := new(big.Int).SetString(str, 2)
			bn.Add(bn, big.NewInt(1))
			bn.Mul(bn, big.NewInt(-1))
			heads[0] = bn
		}

		s = fmt.Sprintf("%0*b", 256, heads[1])
		if string(s[0]) == "1" {
			str := flipStr(s)
			bn, _ := new(big.Int).SetString(str, 2)
			bn.Add(bn, big.NewInt(1))
			bn.Mul(bn, big.NewInt(-1))
			heads[1] = bn
		}

		res := f(heads[0], heads[1])
		if res.Cmp(big.NewInt(0)) == -1 {
			s = fmt.Sprintf("%0*b", 256, res.Mul(res, big.NewInt(-1)))
			str := flipStr(s)
			bn, _ := new(big.Int).SetString(str, 2)
			bn.Add(bn, big.NewInt(1))
			res = bn
		}
		return res
	}
	return f(heads[0], heads[1])
}

func (s *evmstack) getHeads() []*big.Int {
	heads := []*big.Int{s.s[0], s.s[1]}
	s.s = s.s[2:]
	return heads
}

func evm(code []byte) []*big.Int {
	var stack *evmstack = &evmstack{}
	pc := 0

LOOP:
	for pc < len(code) {
		opcode := code[pc]

		switch opcode {
		case 00: // STOP
			break LOOP
		case 0x60: // PUSH1
			pb := 1
			item := fmt.Sprintf("%x", code[pc+1:pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.s = append([]*big.Int{bn}, stack.s...)
			pc += pb
		case 0x50: // POP
			stack.s = stack.s[1:]
		case 0x01: // ADD
			res := stack.oprHeads(new(big.Int).Add, false)
			res.Mod(res, max.uint256Max)
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x7f: // PUSH32
			pb := 32
			item := fmt.Sprintf("%x", code[pc+1:pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.s = append([]*big.Int{bn}, stack.s...)
			pc += pb
		case 0x02: // MUL
			res := stack.oprHeads(new(big.Int).Mul, false)
			res.Mod(res, max.uint256Max)
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x03: // SUB
			res := stack.oprHeads(new(big.Int).Sub, false)
			res.Mod(res, max.uint256Max)
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x04, 0x05: // DIV, SDIV
			var res *big.Int
			if stack.s[1].String() == "0" {
				stack.s = stack.s[2:]
				res = big.NewInt(0)
			} else {
				if opcode == 0x04 {
					res = stack.oprHeads(new(big.Int).Div, false)
				} else {
					res = stack.oprHeads(new(big.Int).Div, true)
				}
			}
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x06, 0x07: // MOD, SMOD
			var res *big.Int
			if stack.s[1].String() == "0" {
				stack.s = stack.s[2:]
				res = big.NewInt(0)
			} else {
				if opcode == 0x06 {
					res = stack.oprHeads(new(big.Int).Rem, false)
				} else {
					res = stack.oprHeads(new(big.Int).Rem, true)
				}
			}
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x10, 0x11: // LT, GT
			heads := stack.getHeads()
			var bn *big.Int
			pb := 1

			cmp := 1
			if opcode == 0x11 {
				cmp = -1
			}

			if heads[1].Cmp(heads[0]) == cmp {
				bn = big.NewInt(1)
			} else {
				bn = big.NewInt(0)
			}
			stack.s = append([]*big.Int{bn}, stack.s...)
			pc += pb
		}
		pc++
	}
	return stack.s
}

func main() {
	content, err := os.ReadFile("./evm.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []TestCase
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during json.Unmarshal(): ", err)
	}

	for index, test := range payload {
		fmt.Printf("Test #%v of %v: %v\n", index+1, len(payload), test.Name)

		bin, err := hex.DecodeString(test.Code.Bin)
		if err != nil {
			log.Fatal("Error during hex.DecodeString(): ", err)
		}

		var expectedStack []*big.Int
		for _, s := range test.Expect.Stack {
			i, ok := new(big.Int).SetString(s, 0)
			if !ok {
				log.Fatal("Error during big.Int.SetString(): ", err)
			}
			expectedStack = append(expectedStack, i)
		}

		// Note: as the test cases get more complex, you'll need to modify this
		// to pass down more arguments to the evm function and return more than
		// just the stack.
		stack := evm(bin)

		match := len(stack) == len(expectedStack)
		if match {
			for i, s := range stack {
				match = match && (s.Cmp(expectedStack[i]) == 0)
			}
		}

		if !match {
			fmt.Printf("Instructions: \n%v\n", test.Code.Asm)
			fmt.Printf("Expected: %v\n", toStrings(expectedStack))
			fmt.Printf("Got: %v\n\n", toStrings(stack))
			fmt.Printf("Progress: %v/%v\n\n", index, len(payload))
			log.Fatal("Stack mismatch")
		}
	}
}
