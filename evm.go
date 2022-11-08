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

	SWITCH:
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
		case 0x61: // PUSH2
			pb := 2
			item := fmt.Sprintf("%x", code[pc+1:pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.s = append([]*big.Int{bn}, stack.s...)
			pc += pb
		case 0x62: // PUSH3
			pb := 3
			item := fmt.Sprintf("%x", code[pc+1:pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.s = append([]*big.Int{bn}, stack.s...)
			pc += pb
		case 0x7f: // PUSH32
			pb := 32
			item := fmt.Sprintf("%x", code[pc+1:pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.s = append([]*big.Int{bn}, stack.s...)
			pc += pb
		case 0x50: // POP
			_ = stack.getHeads(1)[0]
		case 0x01: // ADD
			res := stack.oprHeads(new(big.Int).Add, false)
			res.Mod(res, max.uint256Max)
			stack.s = append([]*big.Int{res}, stack.s...)
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
			heads := stack.getHeads(2)
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
		case 0x12, 0x13: // SLT, SGT
			heads := stack.getHeads(2)
			heads = twosComps(heads)
			var bn *big.Int

			cmp := 1
			if opcode == 0x13 {
				cmp = -1
			}

			if heads[1].Cmp(heads[0]) == cmp {
				bn = big.NewInt(1)
			} else {
				bn = big.NewInt(0)
			}
			stack.s = append([]*big.Int{bn}, stack.s...)
		case 0x14: // EQ
			heads := stack.getHeads(2)
			var bn *big.Int
			if heads[1].Cmp(heads[0]) == 0 {
				bn = big.NewInt(1)
			} else {
				bn = big.NewInt(0)
			}
			stack.s = append([]*big.Int{bn}, stack.s...)
		case 0x15: // ISZERO
			head := stack.getHeads(1)[0]
			var bn *big.Int
			bn = big.NewInt(0)

			if head.Cmp(big.NewInt(0)) == 0 {
				bn = big.NewInt(1)
			}
			stack.s = append([]*big.Int{bn}, stack.s...)
		case 0x16: // AND
			res := stack.oprHeads(new(big.Int).And, false)
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x17: // OR
			res := stack.oprHeads(new(big.Int).Or, false)
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x18: // XOR
			res := stack.oprHeads(new(big.Int).Xor, false)
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x19: // NOT
			res := stack.oprHead(new(big.Int).Not, true)
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x1a: // BYTE
			heads := stack.getHeads(2)
			r := new(big.Int).Sub(big.NewInt(248), new(big.Int).Mul(heads[0], big.NewInt(8)))
			res := new(big.Int).Rsh(heads[1], uint(r.Int64()))
			bn := byteToBn("ff")
			res = new(big.Int).And(res, bn)
			stack.s = append([]*big.Int{res}, stack.s...)
		case 0x80: // DUP1
			head := stack.getHeads(1)[0]
			stack.s = append([]*big.Int{head, head}, stack.s...)
		case 0x82: // DUP3
			heads := stack.getHeads(3)
			stack.s = append([]*big.Int{heads[2], heads[0], heads[1], heads[2]}, stack.s...)
		case 0x90: // SWAP1
			heads := stack.getHeads(2)
			stack.s = append([]*big.Int{heads[1], heads[0]}, stack.s...)
		case 0x92: // SWAP3
			heads := stack.getHeads(4)
			stack.s = append([]*big.Int{heads[3], heads[1], heads[2], heads[0]}, stack.s...)
		case 0xfe: // INVALID
			break SWITCH
		case 0x58: // PC
			bn := big.NewInt(int64(pc))
			stack.s = append([]*big.Int{bn}, stack.s...)
		}
		pc++
	}
	return stack.s
}
