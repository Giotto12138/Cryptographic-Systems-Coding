package main

import (
	"crypto/rand"
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

func main() {

	var one *big.Int = big.NewInt(1)

	args := os.Args

	// check the arguments
	if len(args) != 4 {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: elg-encrypt <message text as a string with quotes> <filename of public key> <filename of ciphertext>")
		return
	}
	AlicePub := args[1]
	BobPub := args[2]

	// read the file storing public key of Alice
	text, err := os.Open(AlicePub)
	if err != nil {
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	pubStr, _ := ioutil.ReadAll(text)

	// get p, g, g_a from the public key of Alice
	p, g, g_a := pubAnalyze(string(pubStr))

	// generate secret number b
	var b *big.Int = new(big.Int)
	var pMinusOne *big.Int = new(big.Int)
	pMinusOne.Sub(p, one)                   // get p-1
	b, _ = rand.Int(rand.Reader, pMinusOne) // get a random b from 1 to p-1

	// get g^b mod p
	var g_b *big.Int = new(big.Int)
	g_b.Exp(g, b, p)

	pubKey := "( " + g_b.String() + " )"

	var f1 *os.File
	var err1 error

	// write public key into the first file
	if checkFileIsExist(BobPub) { // if the file exists, open it
		f1, err1 = os.OpenFile(BobPub, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	} else {
		f1, err1 = os.Create(BobPub) // if the file doesn't exists, create it
	}

	check(err1)
	_, err1 = io.WriteString(f1, pubKey)

	// get g^{ab} mod p
	var g_ab *big.Int = new(big.Int)
	g_ab.Exp(g_a, b, p)

	fmt.Println("the shared secret key g^{ab} mod p")
	fmt.Println(g_ab)

}
