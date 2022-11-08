package evm

import (
	"evm-from-scratch-go/domain"
	"evm-from-scratch-go/utils"
	"fmt"
	"math/big"
)

func Evm(code []byte) []*big.Int {
	var stack *EvmStack = &EvmStack{}
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
			stack.S = append([]*big.Int{bn}, stack.S...)
			pc += pb
		case 0x61: // PUSH2
			pb := 2
			item := fmt.Sprintf("%x", code[pc+1:pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.S = append([]*big.Int{bn}, stack.S...)
			pc += pb
		case 0x62: // PUSH3
			pb := 3
			item := fmt.Sprintf("%x", code[pc+1:pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.S = append([]*big.Int{bn}, stack.S...)
			pc += pb
		case 0x7f: // PUSH32
			pb := 32
			item := fmt.Sprintf("%x", code[pc+1:pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.S = append([]*big.Int{bn}, stack.S...)
			pc += pb
		case 0x50: // POP
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
		case 0x01: // ADD
			res := stack.oprHeads(new(big.Int).Add, false)
			if res == nil {
				break LOOP
			}
			res.Mod(res, domain.Max.Uint256Max)
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x02: // MUL
			res := stack.oprHeads(new(big.Int).Mul, false)
			if res == nil {
				break LOOP
			}
			res.Mod(res, domain.Max.Uint256Max)
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x03: // SUB
			res := stack.oprHeads(new(big.Int).Sub, false)
			if res == nil {
				break LOOP
			}
			res.Mod(res, domain.Max.Uint256Max)
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x04, 0x05: // DIV, SDIV
			var res *big.Int
			if stack.S[1].String() == "0" {
				stack.S = stack.S[2:]
				res = big.NewInt(0)
			} else {
				if opcode == 0x04 {
					res = stack.oprHeads(new(big.Int).Div, false)
				} else {
					res = stack.oprHeads(new(big.Int).Div, true)
				}
				if res == nil {
					break LOOP
				}
			}
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x06, 0x07: // MOD, SMOD
			var res *big.Int
			if stack.S[1].String() == "0" {
				stack.S = stack.S[2:]
				res = big.NewInt(0)
			} else {
				if opcode == 0x06 {
					res = stack.oprHeads(new(big.Int).Rem, false)
				} else {
					res = stack.oprHeads(new(big.Int).Rem, true)
				}
				if res == nil {
					break LOOP
				}
			}
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x10, 0x11: // LT, GT
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

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
			stack.S = append([]*big.Int{bn}, stack.S...)
		case 0x12, 0x13: // SLT, SGT
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			heads = utils.TwosComps(heads)
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
			stack.S = append([]*big.Int{bn}, stack.S...)
		case 0x14: // EQ
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			var bn *big.Int
			if heads[1].Cmp(heads[0]) == 0 {
				bn = big.NewInt(1)
			} else {
				bn = big.NewInt(0)
			}
			stack.S = append([]*big.Int{bn}, stack.S...)
		case 0x15: // ISZERO
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
			var bn *big.Int
			bn = big.NewInt(0)

			if head.Cmp(big.NewInt(0)) == 0 {
				bn = big.NewInt(1)
			}
			stack.S = append([]*big.Int{bn}, stack.S...)
		case 0x16: // AND
			res := stack.oprHeads(new(big.Int).And, false)
			if res == nil {
				break LOOP
			}
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x17: // OR
			res := stack.oprHeads(new(big.Int).Or, false)
			if res == nil {
				break LOOP
			}
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x18: // XOR
			res := stack.oprHeads(new(big.Int).Xor, false)
			if res == nil {
				break LOOP
			}
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x19: // NOT
			res := stack.oprHead(new(big.Int).Not, true)
			if res == nil {
				break LOOP
			}
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x1a: // BYTE
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			r := new(big.Int).Sub(big.NewInt(248), new(big.Int).Mul(heads[0], big.NewInt(8)))
			res := new(big.Int).Rsh(heads[1], uint(r.Int64()))
			bn := utils.ByteToBn("ff")
			res = new(big.Int).And(res, bn)
			stack.S = append([]*big.Int{res}, stack.S...)
		case 0x80: // DUP1
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
			stack.S = append([]*big.Int{head, head}, stack.S...)
		case 0x82: // DUP3
			heads := stack.getHeads(3)
			if heads == nil {
				break LOOP
			}
			stack.S = append([]*big.Int{heads[2], heads[0], heads[1], heads[2]}, stack.S...)
		case 0x90: // SWAP1
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}
			stack.S = append([]*big.Int{heads[1], heads[0]}, stack.S...)
		case 0x92: // SWAP3
			heads := stack.getHeads(4)
			if heads == nil {
				break LOOP
			}
			stack.S = append([]*big.Int{heads[3], heads[1], heads[2], heads[0]}, stack.S...)
		case 0xfe: // INVALID
			break SWITCH
		case 0x58: // PC
			bn := big.NewInt(int64(pc))
			stack.S = append([]*big.Int{bn}, stack.S...)
		case 0x5a: // GAS
			bn := utils.ByteToBn("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
			stack.S = append([]*big.Int{bn}, stack.S...)
		case 0x56: // JUMP
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
			pc = int(head.Int64())
			if fmt.Sprintf("%x", code[pc]) != "5b" {
				break LOOP
			}
		case 0x57: // JUMP1
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}
			if heads[1].Cmp(big.NewInt(0)) != 0 {
				if fmt.Sprintf("%x", code[int(heads[0].Int64())]) == "5b" {
					pc = int(heads[0].Int64())
				} else {
					break LOOP
				}
			}
		}
		pc++
	}
	return stack.S
}
