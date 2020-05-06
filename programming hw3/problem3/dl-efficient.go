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

func efficient(p, g, g_a *big.Int) {
	var zero *big.Int = big.NewInt(0)
	var one *big.Int = big.NewInt(1)
	//var two *big.Int = big.NewInt(2)
	var x *big.Int = new(big.Int)
	var temp *big.Int = new(big.Int)
	var m *big.Int = new(big.Int)
	var j *big.Int = new(big.Int)
	var i *big.Int = new(big.Int)
	var g_j *big.Int = new(big.Int)
	var g_aj *big.Int = new(big.Int)

	// x=im-j, m =ceil(sqrt(p)), if g^{im} == g^a * g^j, then we find x
	m.Sqrt(p)
	temp.Mul(m, m)
	// m^2 should not be less than p
	for true {
		if temp.Cmp(p) >= 0 {
			break
		}
		if temp.Cmp(p) < 0 {
			m.Add(m, one)
			temp.Mul(m, m)
		}
	}
	// fmt.Println("p:  ", p)
	// fmt.Println("temp:  ", temp)

	// a map to store all the g_aj
	dict := make(map[string]*big.Int)

	// iterate j, beginning from 0 to m
	for j.Set(zero); j.Cmp(m) < 0; j.Add(j, one) {
		g_j.Exp(g, j, p)
		g_aj.Mul(g_a, g_j)
		g_aj.Exp(g_aj, one, p)
		//fmt.Println("g_aj: ", g_aj)
		gStr := g_aj.String()
		dict[gStr] = j
	}

	// iterate i, beginning from 0 to m
	for i.Set(zero); i.Cmp(m) < 0; i.Add(i, one) {

		temp.Mul(i, m)
		temp.Exp(g, temp, p)
		//fmt.Println("g_im:  ", temp)

		if element, ok := dict[temp.String()]; ok {
			// calculate x=im-j
			fmt.Println("final")
			fmt.Println("i:  ", i)
			fmt.Println("j:  ", element)
			temp.Mul(i, m)
			x.Sub(temp, element)
			fmt.Println("x is: ", x)
			break
		}
	}

	if i.Cmp(m) == 0 {
		fmt.Println("cannot find x")
	}

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

	// Use Baby steps giant steps algorithm to find a more quickly
	efficient(p, g, g_a)

}
