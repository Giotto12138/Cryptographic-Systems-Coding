package main

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

func xor(block1 []byte, block2 []byte) []byte {
	if len(block1) != len(block2) {
		fmt.Println("Please input same length!")
		os.Exit(1)
	}
	result := make([]byte, len(block1))
	for i := 0; i < len(block1); i++ {
		result[i] = block1[i] ^ block2[i]
	}
	return result
}

func hmac(key []byte, plaintext []byte) [32]byte {
	k_derived := make([]byte, 64)
	ipad := make([]byte, 64)
	opad := make([]byte, 64)
	for i := 0; i < 64; i++ {
		ipad[i] = 0x36
		opad[i] = 0x5c
	}
	if len(key) > 64 {
		hash := sha256.Sum256(key)
		copy(k_derived, hash[0:])
	} else {
		copy(k_derived, key[0:])
	}

	xored_ipad := make([]byte, 64)
	xored_opad := make([]byte, 64)
	xored_opad = xor(opad, k_derived)
	xored_ipad = xor(ipad, k_derived)
	for i := 0; i < len(plaintext); i++ {
		xored_ipad = append(xored_ipad, plaintext[i])
	}
	hash_ipad := sha256.Sum256(xored_ipad)
	final_hash := xored_opad
	for i := 0; i < len(hash_ipad); i++ {
		final_hash = append(final_hash, hash_ipad[i])
	}
	h_mac := sha256.Sum256(final_hash)
	return h_mac
}

func encrypt(plaintext []byte, key_enc []byte, key_mac []byte, outputfile string) {

	c, err := aes.NewCipher(key_enc)
	if err != nil {
		fmt.Println(err)
	}
	mac := hmac(key_mac, plaintext)
	for i := 0; i < len(mac); i++ {
		plaintext = append(plaintext, mac[i])
	}

	block_num := len(plaintext) / 16
	remain := len(plaintext) % 16
	//padding
	block_num++
	remain = 16 - remain
	padding := byte(remain)
	for i := 0; i < remain; i++ {
		plaintext = append(plaintext, padding)
	}
	//Generate IV
	iv := make([]byte, 16)
	_, iv_err := rand.Read(iv)
	if iv_err != nil {
		fmt.Println("Error when creating IV")
		os.Exit(1)
	}

	xored_cipher := make([]byte, 16)
	block_cipher := make([]byte, 16)
	ciphertext := make([]byte, 16)
	//encrypt first block
	xored_cipher = xor(plaintext[0:16], iv)
	c.Encrypt(block_cipher, xored_cipher)
	for i := 0; i < len(block_cipher); i++ {
		ciphertext[i] = block_cipher[i]
	}

	for i := 0; i < block_num; i++ {
		for j := 0; j < len(block_cipher); j++ {
			iv[j] = block_cipher[j]
		}
		xored_cipher = xor(plaintext[16*i:16*(i+1)], iv)
		c.Encrypt(block_cipher, xored_cipher)
		for k := 0; k < len(block_cipher); k++ {
			ciphertext = append(ciphertext, block_cipher[k])
		}

	}

	err_write := ioutil.WriteFile(outputfile, ciphertext, 0777)
	if err_write != nil {
		fmt.Println("Can not write to file!")
	}
	fmt.Println(ciphertext)
	return

}

func decrypt(ciphertext []byte, iv []byte, key_enc []byte, key_mac []byte, outputfile string) {
	c, err := aes.NewCipher(key_enc)
	if err != nil {
		fmt.Println(err)
	}

	if len(ciphertext)%16 != 0 {
		fmt.Println("Invalid CipherText")
		os.Exit(1)
	}
	block_num := len(ciphertext) / 16
	text_with_iv := make([]byte, 16)
	text_with_padding := make([]byte, 16)
	c.Decrypt(text_with_iv, ciphertext[0:16])
	text_with_iv = xor(text_with_iv, iv)
	copy(text_with_padding, text_with_iv)
	for i := 1; i < block_num; i++ {
		c.Decrypt(text_with_iv, ciphertext[i*16:(i+1)*16])
		text_with_iv = xor(text_with_iv, iv)
		for j := 1; j < len(text_with_iv); j++ {
			text_with_padding = append(text_with_padding, text_with_iv[j])
		}
	}
	//strip padding
	padding := text_with_padding[len(text_with_padding)-1]
	for i := (len(text_with_padding) - int(padding)); i < (len(text_with_padding)); i++ {
		//if text_with_padding[i]!=padding{
		//	fmt.Println("Invalid padding")
		//	os.Exit(1)
		//}
	}
	text_with_mac := make([]byte, len(text_with_padding)-int(padding))
	copy(text_with_mac, text_with_padding[0:len(text_with_padding)-int(padding)])
	mac := make([]byte, 32)
	copy(mac, text_with_mac[len(text_with_mac)-32:])
	plaintext := make([]byte, len(text_with_mac)-32)
	copy(plaintext, text_with_mac[0:len(text_with_mac)-32])
	fmt.Println(string(plaintext))
	err_write := ioutil.WriteFile(outputfile, plaintext, 0777)
	if err_write != nil {
		fmt.Println("Can not write to file!")
	}
	return

}

func main() {
	if len(os.Args) != 9 {
		fmt.Println("Invalid input argument")
		fmt.Println("Expected input is: encrypt-auth <mode> -k <32-byte key in hexadecimal> -i <input file> -o <outputfile>")
		os.Exit(1)
	}

	if len(os.Args[4]) != 64 {
		fmt.Println("key length error")
		os.Exit(1)
	}

	key := os.Args[4]
	enc_str := key[0:32]
	mac_str := key[32:64]
	key_enc, _ := hex.DecodeString(enc_str)
	key_mac, _ := hex.DecodeString(mac_str)
	if os.Args[2] == "encrypt" {
		plaintext, _ := ioutil.ReadFile(os.Args[6])
		outputfile := os.Args[8]
		//fmt.Println(plaintext)
		encrypt(plaintext, key_enc, key_mac, outputfile)
		return
	} else if os.Args[2] == "decrypt" {
		rawdata, _ := ioutil.ReadFile(os.Args[6])
		iv := make([]byte, 16)
		iv = rawdata[0:16]
		ciphertext := make([]byte, len(rawdata)-16)
		ciphertext = rawdata[16:]
		outputfile := os.Args[8]
		fmt.Println("ciphertext:", ciphertext)
		decrypt(ciphertext, iv, key_enc, key_mac, outputfile)
		return
	}
}
