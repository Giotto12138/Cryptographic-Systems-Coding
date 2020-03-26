package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func attack(ciphertext []byte) {

	plaintext := make([]byte, len(ciphertext))
	copy(plaintext, ciphertext)

	// calculate how many blocks
	num := len(ciphertext) / 16

	// starting from the end, process every two blocks for one time
	for i := num - 1; i > 0; i-- {

		// move two blocks to the end
		copy(ciphertext[len(ciphertext)-32:len(ciphertext)], plaintext[i*16-16:i*16+16])

		// c1 is the first block of ciphertext of two chosen blocks
		c1 := make([]byte, 16)
		copy(c1, ciphertext[len(ciphertext)-32:len(ciphertext)-16])

		// intermediate is the decrypted block of the second block of two chosen blocks
		intermediate := make([]byte, 16)

		// c_1 is to store modified c1
		c_1 := ciphertext[len(ciphertext)-32 : len(ciphertext)-16]
		_, c1_err := rand.Read(c_1)
		if c1_err != nil {
			fmt.Print(c1_err)
		}

		// try every byte of the block
		for i := 15; i >= 0; i-- {
			pad := byte(16 - i)

			// modify ensure byte of c_1 for next try
			for j := i + 1; j < 16; j++ {
				c_1[j] = pad ^ intermediate[j]
			}

			// try every possible value to guess
			for k := 0x00; k < 0x100; {

				// set guessed byte from 0x00 to 0x100 to find a right one
				c_1[i] = byte(k)

				var output []byte
				output = make([]byte, hex.EncodedLen(len(ciphertext)))
				hex.Encode(output, ciphertext)

				// write output to a file for decrypt-test to decrypt
				writeFile_err := ioutil.WriteFile("attack.txt", output, 0644)
				if writeFile_err != nil {
					panic(writeFile_err)
				}

				// run decypt-test program, and get output to do attacking
				out, test_err := exec.Command("./decrypt-test", "-i", "attack.txt").CombinedOutput()
				if test_err != nil {
					fmt.Print(test_err)
				}
				// if padding is right, get the right byte
				if !strings.Contains(string(out), "INVALID PADDING") {
					break
				}
				k++
			}

			// get the intermediate state
			intermediate[i] = pad ^ c_1[i]
		}

		block := make([]byte, 16)
		// xor intermediate state with c1 to get the plaintext and store in block
		for i := 0; i < 16; i++ {
			block[i] = intermediate[i] ^ c1[i]
		}

		// copy every block to plaintext
		copy(plaintext[i*16:i*16+16], block)

	}

	// output Just the plaintext (no padding, no IV)
	fmt.Println()
	fmt.Print("plaintext:  ")
	fmt.Println(plaintext[16:])

	// write real plaintext into a file
	write_err := ioutil.WriteFile("c.txt", plaintext[16:], 0644)
	if write_err != nil {
		panic(write_err)
	}
}

func main() {

	args := os.Args

	if len(os.Args) != 3 {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: decrypt-attack -i <ciphertext file>")
		return
	}
	// get the ciphertext
	ciphertext, readfile_err := ioutil.ReadFile(args[2])
	if readfile_err != nil {
		fmt.Print(readfile_err)
	}

	attack(ciphertext)

	// output Just the plaintext (no padding, no IV)
}
