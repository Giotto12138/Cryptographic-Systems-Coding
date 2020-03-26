package main

import (
	"bytes"
	"crypto/aes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
)

func hashMac(key_mac []byte, plaintext []byte) [32]byte {

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

func decrypt(input []byte, key_enc []byte, key_mac []byte) {

	// initiate aes block
	block, aes_err := aes.NewCipher(key_enc)
	if aes_err != nil {
		panic(aes_err)
	}

	// get the initial iv
	iv := make([]byte, 16)
	copy(iv, input[:16])
	// get the ciphertext
	ciphertext := make([]byte, len(input)-16)
	copy(ciphertext, input[16:])

	plaintext := make([]byte, len(ciphertext))
	//decrypt every block
	for index := 0; index < len(ciphertext); index = index + 16 {

		kz := make([]byte, 16)
		// decrypt every block of ciphertext and store in kz
		block.Decrypt(kz, ciphertext[index:index+16])
		// xor kz and iv
		j := 0
		for i := index; i < index+16; i++ {
			plaintext[i] = kz[j] ^ iv[j]
			j = j + 1
		}
		// previous ciphertext is a new iv
		iv = ciphertext[index : index+16]
	}

	// check padding
	padded := plaintext[len(plaintext)-1]
	pad_check := make([]byte, int(padded))
	// get all the padded value and store them in pad_check
	copy(pad_check, plaintext[len(plaintext)-int(padded):len(plaintext)])

	for i := 0; i < int(padded); i++ {
		if pad_check[i] != padded {
			fmt.Println("INVALID PADDING")
			os.Exit(1)
		}
	}

	// get the real plaintext by strip padding and 32 byte tag
	real_plaintext := make([]byte, len(plaintext)-int(padded)-32)
	copy(real_plaintext, plaintext[:len(plaintext)-int(padded)-32])

	// calculate the tag
	tag := hashMac(key_mac, real_plaintext)
	tag_check := plaintext[len(plaintext)-int(padded)-32 : len(plaintext)-int(padded)]

	for i := 0; i < len(tag); i++ {
		if tag[i] != tag_check[i] {
			fmt.Println("INVALID MAC")
			os.Exit(1)
		}
	}

	// // write real plaintext into a file
	// write_err := ioutil.WriteFile(output, real_plaintext, 0644)
	// if write_err != nil {
	// 	panic(write_err)
	// }

	fmt.Println("“SUCCESS” ")
}

func main() {

	args := os.Args

	// check the arguments
	if len(args) != 3 || args[1] != "-i" {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: decrypt-test -i <ciphertext file>")
		return
	}

	// first 16 bytes are encryption key, last 16 bytes are key to mac
	key := "abcdeasagdguyesfabsdeugjkdnxlqos"
	key_enc := []byte(key[0:16])
	key_mac := []byte(key[16:32])

	ciphertext, readfile_err := ioutil.ReadFile(args[2])
	if readfile_err != nil {
		fmt.Print(readfile_err)
	}

	decrypt(ciphertext, key_enc, key_mac)

}
