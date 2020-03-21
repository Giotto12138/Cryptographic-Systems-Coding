package main

import (
	"os"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"bytes"
	"io/ioutil"
)


func encrypt(plaintext []byte, key_enc []byte, key_mac []byte, output string) {

	c, err := aes.NewCipher(key_enc)
	if err != nil {
		fmt.Println(err)
	}
	mac := hmac (key_mac, plaintext)
	for i:= 0; i < len(mac); i++ {
		plaintext = append(plaintext,mac[i])
	}

	block_num := len(plaintext)/16
	remain := len(plaintext)%16
	//padding
	block_num++
	remain = 16-remain
	padding := byte(remain)
	for i := 0; i < remain; i++ {
		plaintext = append(plaintext, padding)
	}
	//Generate IV
	iv := make([]byte,16)
	_,iv_err := rand.Read(iv)
	if iv_err != nil {
		fmt.Println("Error when creating IV")
		os.Exit(1)
	}

	xored_cipher := make([]byte,16)
	block_cipher := make([]byte,16)
	ciphertext := make([]byte,16)
	//encrypt first block
	xored_cipher  = xor(plaintext[0:16],iv)
	c.Encrypt(block_cipher,xored_cipher)
	for i:=0; i<len(block_cipher); i++{
		ciphertext[i] = block_cipher[i]
	}

	for i:=0; i < block_num; i++{
		for j:=0; j<len(block_cipher); j++{
			iv[j] = block_cipher[j]
		}
		xored_cipher  = xor(plaintext[16*i:16*(i+1)],iv)
		c.Encrypt(block_cipher,xored_cipher)
		for k:=0; k<len(block_cipher); k++{
			ciphertext = append(ciphertext, block_cipher[k])
		}

	}

	err_write := ioutil.WriteFile(outputfile, ciphertext, 0777)
	if err_write !=nil{
		fmt.Println("Can not write to file!")
	}
	fmt.Println(ciphertext)
	return
}


func main() {

	args := os.Args

	if len(args) != 9 {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: encrypt-auth <mode> -k <32-byte key in hexadecimal> -i <input file> -o <outputfile>")
		return
	}

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
	else if args[2] == "decrypt" {
		fmt.Print("decrypt")
	}
}