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

func bruteForce(p, g, g_a *big.Int) {
	var zero *big.Int = big.NewInt(0)
	var one *big.Int = big.NewInt(1)
	var pMinusOne = new(big.Int).Sub(p, one)
	var x *big.Int = new(big.Int)
	var temp *big.Int = new(big.Int)

	// guess the x, beginning from p-1 to 1
	for x.Set(pMinusOne); x.Cmp(zero) > 0; x.Sub(x, one) {
		temp.Exp(g, x, p)
		//fmt.Println("temp:", temp)

		if temp.Cmp(g_a) == 0 {
			break
		}
	}

	fmt.Println("x is: ", x)

}

func main() {

	args := os.Args

	// check the arguments
	if len(args) != 2 {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: dl-brute <filename for inputs>")
		return
	}
	AlicePub := args[1]

	// read the file storing public key
	text, err := os.Open(AlicePub)
	if err != nil {
		fmt.Println("Input File Error")
		os.Exit(1)
	}
	pubStr, _ := ioutil.ReadAll(text)

	// get p, g, g_a from the public key
	p, g, g_a := pubAnalyze(string(pubStr))

	bruteForce(p, g, g_a)

}
