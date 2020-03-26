package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func attack(ciphertext []byte) {

	// get iv from ciphertext
	iv := ciphertext[:16]
	ciphertext = ciphertext[16:]

	// use decrypt-test to check padding
	for i := 0; i < len(ciphertext); i = i + 16 {
		// decrypt every block(16 bytes) of ciphertext
		block := ciphertext[i : i+16]

		// modify one byte of iv every time
		for j := 0; j < 16; j++ {
			for k := 16 - j - 1; k < 16; k++ {
				iv[k] = byte(j)
			}
		}

		// update iv,
	}

	// output Just the plaintext (no padding, no IV)
	fmt.Println(plaintext)
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
