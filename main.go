/**
 * EVM From Scratch Program
 * implemented in Go by restructring the draft
 * from "github.com/w1nt3r-eth/evm-from-scratch"
 * check README for the program details
 * @author "github.com/ulerdogan"
 * to run ```go run .```
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
