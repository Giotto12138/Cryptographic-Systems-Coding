package main

import (
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

// analyze the public key from Alice
func pubAnalyze(pub string) *big.Int {

	var g_b big.Int

	// delete spaces in the string
	pub = strings.Replace(pub, " ", "", -1)

	temp := strings.Split(pub, "(")
	temp = strings.Split(temp[1], ")")

	// convert string to *big.Int
	g_b.SetString(temp[0], 10)

	return &g_b
}

// analyze the private key from Alice
func priAnalyze(pri string) (*big.Int, *big.Int) {

	var a, p big.Int

	// delete spaces in the string
	pri = strings.Replace(pri, " ", "", -1)
	// split the string with comma, get three parts
	strSplit := strings.Split(pri, ",")
	// get p from the first part
	pStr := strings.Split(strSplit[0], "(")
	// get g_a from the second part
	aStr := strings.Split(strSplit[2], ")")

	// convert string to *big.Int
	p.SetString(pStr[1], 10)
	a.SetString(aStr[0], 10)

	return &a, &p
}

func main() {

	args := os.Args

	// check the arguments
	if len(args) != 3 {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: dh-alice2 <filename of message from Bob> <filename to read secret key>")
		return
	}
	BobPub := args[1]
	AlicePri := args[2]

	// read the file storing public key of Bob
	text1, err := os.Open(BobPub)
	if err != nil {
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	pubStr, _ := ioutil.ReadAll(text1)

	// get g_b from the public key of Bob
	g_b := pubAnalyze(string(pubStr))

	// read the file storing private key of Alice
	text2, err := os.Open(AlicePri)
	if err != nil {
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	priStr, _ := ioutil.ReadAll(text2)

	// get a, p from the public key of Bob
	a, p := priAnalyze(string(priStr))

	// get g^{ab} mod p
	var g_ab *big.Int = new(big.Int)
	g_ab.Exp(g_b, a, p)

	fmt.Println("the shared secret key g^{ab} mod p")
	fmt.Println(g_ab)

}
