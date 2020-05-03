package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

// analyze the public key from Alice
func pubAnalyze(pub string) (*big.Int, *big.Int, *big.Int) {

	var p, g, g_a big.Int

	// delete spaces in the string
	pub = strings.Replace(pub, " ", "", -1)
	// split the string with comma, get three parts
	strSplit := strings.Split(pub, ",")
	// get p from the first part
	pStr := strings.Split(strSplit[0], "(")
	// get g from the second part
	gStr := strSplit[1]
	// get g_a from the second part
	g_aStr := strings.Split(strSplit[2], ")")

	// convert string to *big.Int
	p.SetString(pStr[1], 10)
	g.SetString(gStr, 10)
	g_a.SetString(g_aStr[0], 10)

	return &p, &g, &g_a
}

func hash(g_a, g_b, g_ab *big.Int) []byte {

	temp := g_a.String() + " " + g_b.String() + " " + g_ab.String()
	//fmt.Println("temp:   ", temp)
	res := sha256.Sum256([]byte(temp))
	return res[:]
}

// aes-gcm encryption
func encrypt(plaintext string, key []byte) []byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// generate 16 bits random nonce
	nonce := make([]byte, 16)
	_, _ = rand.Read(nonce)
	if false {
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err.Error())
		}
	}

	//fmt.Printf("nonce: %x\n", nonce)
	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	//fmt.Printf("cipher:%x\n", ciphertext)

	// concatenate nonce with ciphertext
	result := make([]byte, len(nonce)+len(ciphertext))
	result = nonce
	for i := 0; i < len(ciphertext); i++ {
		result = append(result, ciphertext[i])
	}

	return result

}

func main() {

	var one *big.Int = big.NewInt(1)

	args := os.Args

	// check the arguments
	if len(args) != 4 {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: elg-encrypt <message text as a string with quotes> <filename of public key> <filename of ciphertext>")
		return
	}
	plaintext := args[1]
	pubKeyFile := args[2]
	cipherFile := args[3]

	// read the file storing public key
	text, err := os.Open(pubKeyFile)
	if err != nil {
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	pubStr, _ := ioutil.ReadAll(text)

	// get p, g, g_a from the public key
	p, g, g_a := pubAnalyze(string(pubStr))

	// generate secret number b
	var b *big.Int = new(big.Int)
	var pMinusOne *big.Int = new(big.Int)
	pMinusOne.Sub(p, one)                   // get p-1
	b, _ = rand.Int(rand.Reader, pMinusOne) // get a random b from 1 to p-1

	// get g^b mod p
	var g_b *big.Int = new(big.Int)
	g_b.Exp(g, b, p)

	// get g^{ab} mod p as private key
	var g_ab *big.Int = new(big.Int)
	g_ab.Exp(g_a, b, p)

	fmt.Println("the shared secret key g^{ab} mod p")
	fmt.Println(g_ab)

	// generate k = SHA256(g_a||g_b||g_ab).
	k := hash(g_a, g_b, g_ab)
	//fmt.Println(k)

	// encrypt the plaintext with aes-gcm
	ciphertext := encrypt(plaintext, k)
	//fmt.Println(ciphertext)

	// change ciphertext to hexadecimal format
	cipherHex := make([]byte, hex.EncodedLen(len(ciphertext)))
	hex.Encode(cipherHex, ciphertext)
	//fmt.Println(cipherHex)

	var f1 *os.File
	var err1 error

	// write ciphertext into the third file
	if checkFileIsExist(cipherFile) { // if the file exists, open it
		f1, err1 = os.OpenFile(cipherFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	} else {
		f1, err1 = os.Create(cipherFile) // if the file doesn't exists, create it
	}

	check(err1)
	_, err1 = io.WriteString(f1, string(cipherHex))

}
