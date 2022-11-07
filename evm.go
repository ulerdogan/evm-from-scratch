package main

import (
	"fmt"
	"math/big"
)

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
		}
		pc++
	}
	return stack.s
}
