package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * check if a file exists
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// analyze the secret key
func secretAnalyze(secret string) (*big.Int, *big.Int, *big.Int) {

	var p, g, a big.Int

	// delete spaces in the string
	secret = strings.Replace(secret, " ", "", -1)
	// split the string with comma, get three parts
	strSplit := strings.Split(secret, ",")
	// get p from the first part
	pStr := strings.Split(strSplit[0], "(")
	// get g from the second part
	gStr := strSplit[1]
	// get g_a from the second part
	aStr := strings.Split(strSplit[2], ")")

	// convert string to *big.Int
	p.SetString(pStr[1], 10)
	g.SetString(gStr, 10)
	a.SetString(aStr[0], 10)

	return &p, &g, &a
}

// analyze the cipher pair
func cipherAnalyze(pri string) (*big.Int, string) {

	var g_b big.Int

	// delete spaces in the string
	pri = strings.Replace(pri, " ", "", -1)
	// split the string with comma, get three parts
	strSplit := strings.Split(pri, ",")
	// get g_b from the first part
	g_bStr := strings.Split(strSplit[0], "(")
	// get ciphertext from the second part
	cipher := strings.Split(strSplit[1], ")")

	// convert string to *big.Int
	g_b.SetString(g_bStr[1], 10)
	// get the ciphertext
	ciphertext := cipher[0]

	return &g_b, ciphertext
}

func hash(g_a, g_b, g_ab *big.Int) []byte {

	temp := g_a.String() + " " + g_b.String() + " " + g_ab.String()
	//fmt.Println("temp:   ", temp)
	res := sha256.Sum256([]byte(temp))
	return res[:]
}

// aes-gcm decryption
func decrypt(cipherNonce string, key []byte) string {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//fmt.Printf("nonce: %x\n", nonce)
	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		panic(err.Error())
	}

	// get ciphertext and nonce
	nonce := make([]byte, 16)
	ciphertext := make([]byte, len(cipherNonce)-16)
	copy(nonce, cipherNonce[0:16])
	copy(ciphertext, cipherNonce[16:])

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintext)

}

func main() {

	args := os.Args

	// check the arguments
	if len(args) != 3 {
		fmt.Println("Invalid input")
		fmt.Println("elg-decrypt <filename of ciphertext> <filename to read secret key>")
		return
	}
	cipherFile := args[1]
	secretFile := args[2]

	// read the file storing the secret key
	text1, err1 := os.Open(secretFile)
	if err1 != nil {
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	secretStr, _ := ioutil.ReadAll(text1)

	// get p, g, a from the secret key
	p, g, a := secretAnalyze(string(secretStr))

	// read the file storing the ciphertext pair
	text2, err2 := os.Open(cipherFile)
	if err2 != nil {
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	cipherStr, _ := ioutil.ReadAll(text2)

	// get a, p from the secret key
	g_b, ciphertextHex := cipherAnalyze(string(cipherStr))

	// get g^{ab} mod p as private key
	var g_ab *big.Int = new(big.Int)
	g_ab.Exp(g_b, a, p)

	// get g^a mod p
	var g_a *big.Int = new(big.Int)
	g_a.Exp(g, a, p)

	// generate k = SHA256(g_a||g_b||g_ab).
	k := hash(g_a, g_b, g_ab)
	//fmt.Println(k)

	// change ciphertext from hex to decimal
	ciphertext := make([]byte, hex.DecodedLen(len(ciphertextHex)))
	hex.Decode(ciphertext, []byte(ciphertextHex))

	// encrypt the plaintext with aes-gcm
	plaintext := decrypt(string(ciphertext), k)
	fmt.Println(plaintext)

}
