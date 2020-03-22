package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
)

func hmac(key_mac []byte, plaintext []byte) [32]byte {

	// key xor 0x36
	kx := make([]byte, 64)
	for i := 0; i < len(key_mac); i++ {
		kx[i] = key_mac[i] ^ 0x36
	}
	// for rest of byte, fill with 0x36 to get final kx
	for i := len(key_mac); i < 64; i++ {
		kx[i] = 0x36
	}

	// concatenate kx and the plaintext
	var buffer bytes.Buffer
	buffer.Write(kx)
	buffer.Write(plaintext)
	kx_new := buffer.Bytes()

	// get a 32 byte out from kx_new by sha256
	out := sha256.Sum256(kx_new)

	// key xor 0x5C
	ky := make([]byte, 64)
	for i := 0; i < len(key_mac); i++ {
		ky[i] = key_mac[i] ^ 0x5c
	}
	// for rest of byte, fill with 0x5c to get final ky
	for i := len(key_mac); i < 64; i++ {
		ky[i] = 0x5c
	}

	// concatenate ky and the out
	for i := 0; i < len(out); i++ {
		ky = append(ky, out[i])
	}
	ky_new := ky

	// get a 32 byte out from kx_new by sha256
	tag := sha256.Sum256(ky_new)

	return tag
}

func padding(plaintext []byte) []byte {

	n := len(plaintext) % 16

	if n != 0 {
		for i := 0; i < (16 - n); i++ {
			plaintext = append(plaintext, byte(16-n))
		}
	} else {
		for i := 0; i < 16; i++ {
			plaintext = append(plaintext, byte(16))
		}
	}

	return plaintext
}

func encrypt(plaintext []byte, key_enc []byte, key_mac []byte, output string) {

	tag := hmac(key_mac, plaintext)

	// concatenate plaintex and the tag
	for i := 0; i < len(tag); i++ {
		plaintext = append(plaintext, tag[i])
	}

	// padding the plaintext
	plaintext = padding(plaintext)
	//fmt.Print(plaintext)

	// generate random IV
	iv := make([]byte, 16)
	_, iv_err := rand.Read(iv)
	if iv_err != nil {
		fmt.Println("wrong iv")
		os.Exit(1)
	}

}

func main() {

	args := os.Args

	// check the arguments
	if len(args) != 9 {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: encrypt-auth <mode> -k <32-byte key in hexadecimal> -i <input file> -o <outputfile>")
		return
	}

	// first 16 bytes are encryption key, last 16 bytes are key to mac
	key := args[4]
	key_enc := []byte(key[0:16])
	key_mac := []byte(key[16:32])

	if args[2] == "encrypt" {
		plaintext, err := ioutil.ReadFile(args[6])
		if err != nil {
			fmt.Print(err)
		}
		output := args[8]
		encrypt(plaintext, key_enc, key_mac, output)
	}
	if args[2] == "decrypt" {
		fmt.Print("decrypt")
	}
}
