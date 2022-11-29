package evm

import (
	"encoding/hex"
	"evm-from-scratch-go/domain"
	"evm-from-scratch-go/utils"
	"fmt"
	"math/big"
)

func Evm(code []byte, state map[string]domain.AccState, block *domain.BlockInfo, tx *domain.TxData, d interface{}) (result *domain.Result) {
	var stack *EvmStack = &EvmStack{}
	var memory *EvmMemory = &EvmMemory{}
	var logs *domain.Logs = &domain.Logs{}
	result = &domain.Result{Success: true}
	var lastResult *domain.Result = &domain.Result{}
	pc := 0

	var storage map[string]*big.Int = make(map[string]*big.Int)
	if d.(map[string]*big.Int) != nil {
		storage = d.(map[string]*big.Int)
	}

LOOP:
	for pc < len(code) {
		opcode := code[pc]

	SWITCH:
		switch opcode {
		case STOP:
			break LOOP
		case PUSH1, PUSH2, PUSH3, PUSH20, PUSH32:
			pb := int(opcode - 95)
			item := utils.ToHex(code[pc+1 : pc+1+pb])
			bn := utils.HexToBn(item)
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
				if opcode == DIV {
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
				if opcode == MOD {
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
			if opcode == GT {
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
			if opcode == SGT {
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
		case MLOAD:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}
			l := memory.load(int(head.Int64()), 32)
			stack.Stack = append([]*big.Int{l}, stack.Stack...)
		case MSTORE, MSTORE8:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			size := 32
			if opcode == MSTORE8 {
				size = 1
			}

			memory.store(int(heads[0].Int64()), size, heads[1])
		case MSIZE:
			bn := big.NewInt(int64(len(memory.Data)))
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case SHA3:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			m := memory.load(int(heads[0].Int64()), int(heads[1].Int64()))
			bn := utils.Keccak256(m)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case ADDRESS:
			bn := utils.HexToBn(tx.To)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case CALLER:
			bn := utils.HexToBn(tx.From)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case BALANCE:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}

			addr := fmt.Sprintf("0x%s", utils.ToHex(head))
			balance := state[addr].Balance

			bn := utils.HexToBn(balance)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case ORIGIN:
			bn := utils.HexToBn(tx.Origin)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case COINBASE:
			bn := utils.HexToBn(block.Coinbase)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case TIMESTAMP:
			bn := utils.HexToBn(block.Timestamp)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case NUMBER:
			bn := utils.HexToBn(block.Number)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case DIFFICULTY:
			bn := utils.HexToBn(block.Difficulty)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case GASLIMIT:
			bn := utils.HexToBn(block.GasLimit)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case GASPRICE:
			bn := utils.HexToBn(tx.GasPrice)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case CHAINID:
			bn := utils.HexToBn(block.ChainId)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case CALLVALUE:
			bn := utils.HexToBn(tx.Value)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case CALLDATALOAD:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}

			data := tx.Data[int(head.Int64())*2:]
			data = utils.PadRight(data, 64)

			bn := new(big.Int)
			bn.SetString(data, 16)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case CALLDATASIZE:
			size := len(tx.Data)
			bn := big.NewInt(int64(size) / 2)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case CALLDATACOPY:
			heads := stack.getHeads(3)
			if heads == nil {
				break LOOP
			}

			data := tx.Data[heads[1].Int64()*2 : (heads[1].Int64()+heads[2].Int64())*2]
			data = utils.PadRight(data, 64)

			bn := new(big.Int)
			bn.SetString(data, 16)
			memory.store(int(heads[0].Int64()), int(heads[2].Int64()), bn)
		case CODESIZE:
			size := len(code)
			bn := big.NewInt(int64(size))
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case CODECOPY:
			heads := stack.getHeads(3)
			if heads == nil {
				break LOOP
			}

			if int(heads[2].Int64()) > len(code) {
				heads[2] = big.NewInt(int64(len(code)))
			}

			end := int(heads[1].Int64()) + int(heads[2].Int64())
			if end > len(code) {
				end = len(code)
			}

			data := code[int(heads[1].Int64()):end]
			str := utils.ToHex(data)
			str = utils.PadRight(str, int(heads[2].Int64()))

			bn := utils.HexToBn(str)
			memory.store(int(heads[0].Int64()), int(heads[2].Int64()), bn)
		case EXTCODESIZE:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}

			addr := fmt.Sprintf("0x%s", utils.ToHex(head))
			len := len(state[addr].Code.Bin) / 2

			bn := big.NewInt(int64(len))
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case EXTCODECOPY:
			heads := stack.getHeads(4)
			if heads == nil {
				break LOOP
			}

			addr := fmt.Sprintf("0x%s", utils.ToHex(heads[0]))
			extc := state[addr].Code.Bin
			hx, _ := hex.DecodeString(extc)

			if int(heads[3].Int64()) > len(hx) {
				heads[3] = big.NewInt(int64(len(hx)))
			}

			end := int(heads[2].Int64()) + int(heads[3].Int64())
			if end > len(hx) {
				end = len(hx)
			}

			data := hx[int(heads[2].Int64()):end]
			str := utils.ToHex(data)
			str = utils.PadRight(str, int(heads[3].Int64()))

			bn := utils.HexToBn(str)
			memory.store(int(heads[1].Int64()), int(heads[3].Int64()), bn)
		case SELFBALANCE:
			addr := tx.To
			balance := state[addr].Balance

			bn := utils.HexToBn(balance)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case SLOAD:
			head := stack.getHeads(1)[0]
			if head == nil {
				break LOOP
			}

			bn := storage[head.String()]
			if bn == nil {
				bn = big.NewInt(0)
			}
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case SSTORE:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			storage[heads[0].String()] = heads[1]
		case LOG0, LOG1, LOG2, LOG3, LOG4:
			heads := stack.getHeads(int(opcode) - 158)
			if heads == nil {
				break LOOP
			}

			logs.Address = tx.To
			logs.Data = memory.load(int(heads[0].Int64()), int(heads[1].Int64())).String()

			if opcode > 160 {
				topics := heads[2:]
				for i := range topics {
					t := utils.ToAddress(topics[i])
					logs.Topics = append(logs.Topics, t)
				}
			}
		case RETURN:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			bn := memory.load(int(heads[0].Int64()), int(heads[1].Int64()))
			result.Return, result.Success = utils.ToHex(bn), true
		case REVERT:
			heads := stack.getHeads(2)
			if heads == nil {
				break LOOP
			}

			bn := memory.load(int(heads[0].Int64()), int(heads[1].Int64()))
			result.Return, result.Success = utils.ToHex(bn), false
		case CALL:
			heads := stack.getHeads(7)
			if heads == nil {
				break LOOP
			}

			addr := utils.ToAddress(heads[1])
			extc := state[addr].Code.Bin

			newTx := tx
			tx.From = tx.To
			tx.To = addr

			hx, _ := hex.DecodeString(extc)
			res := Evm(hx, state, block, newTx, storage)

			bn := big.NewInt(0)
			if res.Success {
				bn = big.NewInt(1)
			}
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)

			if res.Return != "" {
				m, _ := new(big.Int).SetString(res.Return, 16)
				memory.store(int(heads[5].Int64()), int(heads[6].Int64()), m)
			}
			lastResult.Return = res.Return
		case RETURNDATASIZE:
			size := int64(len(lastResult.Return) / 2)
			bn := big.NewInt(size)
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)
		case RETURNDATACOPY:
			heads := stack.getHeads(3)
			if heads == nil {
				break LOOP
			}

			lr := lastResult.Return[int(heads[1].Int64())*2:]
			bn, _ := new(big.Int).SetString(lr, 16)
			memory.store(int(heads[0].Int64()), int(heads[2].Int64()), bn)
		case DELEGATECALL:
			heads := stack.getHeads(6)
			if heads == nil {
				break LOOP
			}

			extc := state[utils.ToAddress(heads[1])].Code.Bin
			hx, _ := hex.DecodeString(extc)

			newTx := tx
			res := Evm(hx, state, block, newTx, storage)

			bn := big.NewInt(0)
			if res.Success {
				bn = big.NewInt(1)
			}
			stack.Stack = append([]*big.Int{bn}, stack.Stack...)

			if res.Return != "" {
				m, _ := new(big.Int).SetString(res.Return, 16)
				memory.store(int(heads[4].Int64()), int(heads[5].Int64()), m)
			}
			lastResult.Return = res.Return
		}
		pc++
	}
	result.Stack = stack.Stack
	return
}
