package evm

import (
	"evm-from-scratch-go/domain"
	"evm-from-scratch-go/utils"
	"math/big"
)

func Evm(code []byte) []*big.Int {
	var stack *EvmStack = &EvmStack{}
	var memory *EvmMemory = &EvmMemory{}
	pc := 0

LOOP:
	for pc < len(code) {
		opcode := code[pc]

	SWITCH:
		switch opcode {
		case STOP:
			break LOOP
		case PUSH1:
			pb := 1
			item := utils.ToHex(code[pc+1 : pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
			pc += pb
		case PUSH2:
			pb := 2
			item := utils.ToHex(code[pc+1 : pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
			pc += pb
		case PUSH3:
			pb := 3
			item := utils.ToHex(code[pc+1 : pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
			pc += pb
		case PUSH32:
			pb := 32
			item := utils.ToHex(code[pc+1 : pc+1+pb])
			bn := new(big.Int)
			bn.SetString(item, 16)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
			pc += pb
		case POP:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
		case ADD:
			res := stack.oprHeads(new(big.Int).Add, false)
			if res == nil {
				break LOOP
			}
			res.Mod(res, domain.Max.Uint256Max)
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case MUL:
			res := stack.oprHeads(new(big.Int).Mul, false)
			if res == nil {
				break LOOP
			}
			res.Mod(res, domain.Max.Uint256Max)
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case SUB:
			res := stack.oprHeads(new(big.Int).Sub, false)
			if res == nil {
				break LOOP
			}
			res.Mod(res, domain.Max.Uint256Max)
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case DIV, SDIV:
			var res *big.Int
			if stack.Stack[1].String() == "0" {
				stack.Stack = stack.Stack[2:]
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
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case MOD, SMOD:
			var res *big.Int
			if stack.Stack[1].String() == "0" {
				stack.Stack = stack.Stack[2:]
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
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case LT, GT:
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
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case SLT, SGT:
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
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case EQ:
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
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case ISZERO:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
			var bn *big.Int
			bn = big.NewInt(0)

			if head.Cmp(big.NewInt(0)) == 0 {
				bn = big.NewInt(1)
			}
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case AND:
			res := stack.oprHeads(new(big.Int).And, false)
			if res == nil {
				break LOOP
			}
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case OR:
			res := stack.oprHeads(new(big.Int).Or, false)
			if res == nil {
				break LOOP
			}
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case XOR:
			res := stack.oprHeads(new(big.Int).Xor, false)
			if res == nil {
				break LOOP
			}
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case NOT:
			res := stack.oprHead(new(big.Int).Not, true)
			if res == nil {
				break LOOP
			}
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case BYTE:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			r := new(big.Int).Sub(big.NewInt(248), new(big.Int).Mul(heads[0], big.NewInt(8)))
			res := new(big.Int).Rsh(heads[1], uint(r.Int64()))
			bn := utils.ByteToBn("ff")
			res = new(big.Int).And(res, bn)
			stack.Stack = append([]*big.Int{res}, stack.Stack...)
		case DUP1:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
			stack.Stack = append([]*big.Int{head, head}, stack.Stack...)
		case DUP3:
			heads := stack.getHeads(3)
			if heads == nil {
				break LOOP
			}
			stack.Stack = append([]*big.Int{heads[2], heads[0], heads[1], heads[2]}, stack.Stack...)
		case SWAP1:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}
			stack.Stack = append([]*big.Int{heads[1], heads[0]}, stack.Stack...)
		case SWAP3:
			heads := stack.getHeads(4)
			if heads == nil {
				break LOOP
			}
			stack.Stack = append([]*big.Int{heads[3], heads[1], heads[2], heads[0]}, stack.Stack...)
		case INVALID:
			break SWITCH
		case PC:
			bn := big.NewInt(int64(pc))
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case GAS:
			bn := utils.ByteToBn("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case JUMP:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
			pc = int(head.Int64())
			if utils.ToHex(code[pc]) != "5b" {
				break LOOP
			}
		case JUMPI:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}
			if heads[1].Cmp(big.NewInt(0)) != 0 {
				if utils.ToHex(code[int(heads[0].Int64())]) == "5b" {
					pc = int(heads[0].Int64())
				} else {
					break LOOP
				}
			}
		case MLOAD: // MLOAD
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
			l := memory.load(int(head.Int64()))
			stack.Stack = append([]*big.Int{l}, stack.Stack...)
		case MSTORE, MSTORE8:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			size := 32
			if opcode == 0x53 {
				size = 1
			}

			memory.store(int(heads[0].Int64()), size, heads[1])
		case MSIZE:
			bn := big.NewInt(int64(len(memory.Data)))
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		}
		pc++
	}
	return stack.Stack
}
