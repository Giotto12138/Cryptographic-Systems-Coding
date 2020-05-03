package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
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

// generate a 1023 bits prime number q
func getPrime() *big.Int {

	q, err := rand.Prime(rand.Reader, 1023)
	if err != nil {
		fmt.Printf("Can't generate %d-bit prime: %v", 1023, err)
	}
	return q
}

// check a big prime number
func checkPrime(q *big.Int) bool {

	// calculate how many bits the p has
	if q.BitLen() != 1023 {
		fmt.Printf("%v is not %d-bit", q, 1023)
		return false
	}
	// 32 times of Miller Rabin Prime Number test were performed on P. If the method returns true,
	// the probability that P is prime is 1 - (1 / 4) * * 32; otherwise, P is not prime
	if !q.ProbablyPrime(32) {
		fmt.Printf("%v is not prime", q)
		return false
	}

	return true
}

// calculate a  random generator
func generator(p *big.Int, q *big.Int) *big.Int {
	var one *big.Int = big.NewInt(1)
	var two *big.Int = big.NewInt(2)
	//var g *big.Int = big.NewInt(2) // the generator to be checked, begin from 2
	var temp *big.Int = new(big.Int)

	// get a random number from 2 to p-1
	var g *big.Int = new(big.Int)
	var pMinusOne *big.Int = new(big.Int)
	pMinusOne.Sub(p, one)                   // get p-1
	g, _ = rand.Int(rand.Reader, pMinusOne) // get a random a from 1 to p-1

	for true {
		// the formula to check if it is a generator: g{(p-1)/n} mod p
		// n is 2 or q
		flag := 0 // if flag equals 2, it is a generator
		var n *big.Int = new(big.Int)

		for i := 0; i < 2; i++ {

			if i == 0 {
				n.Set(two)
			} else {
				n.Set(q)
			}

			temp.Sub(p, one)     // get p-1
			temp.Div(temp, n)    // get (p-1)/n
			temp.Exp(g, temp, p) // get g{(p-1)/n} mod p

			//fmt.Println("temp: ", temp)

			if temp.Cmp(one) != 0 { // if the result is not 1, it is good
				flag = flag + 1
			}

		}

		if flag == 2 { // this is a generator
			break
		}

		// if not a generator, try a new one
		g, _ = rand.Int(rand.Reader, pMinusOne)
	}

	return g

}

func main() {

	args := os.Args

	// check the arguments
	if len(args) != 3 {
		fmt.Println("Invalid input")
		fmt.Println("The input should be: elg-keygen <filename to store public key> <filename to store secret key>")
		return
	}
	pubFile := args[1]
	priFile := args[2]

	var one *big.Int = big.NewInt(1)
	var two *big.Int = big.NewInt(2)
	var q *big.Int
	var p *big.Int = new(big.Int)

	q = getPrime()

	// check the prime q, if it is not qualified, generate a new q and check again until we have a qualified q
	for true {
		if checkPrime(q) {
			break
		} else {
			q = getPrime()
		}
	}

	// generate a safe prime number p based on q, p = q*2+1
	//temp = q
	p.Mul(q, two)
	p.Add(p, one)

	// fmt.Println("q: ", q)
	// fmt.Println("p: ", p)

	g := generator(p, q)

	// fmt.Println("g: ", g)

	// generate secret number a
	var a *big.Int = new(big.Int)
	var pMinusOne *big.Int = new(big.Int)
	pMinusOne.Sub(p, one)                   // get p-1
	a, _ = rand.Int(rand.Reader, pMinusOne) // get a random a from 1 to p-1

	// get g^a mod p
	var g_a *big.Int = new(big.Int)
	g_a.Exp(g, a, p)

	pubKey := "( " + p.String() + "," + g.String() + "," + g_a.String() + " )"
	priKey := "( " + p.String() + "," + g.String() + "," + a.String() + " )"

	var f1, f2 *os.File
	var err1, err2 error

	// write public key into the first file
	if checkFileIsExist(pubFile) { // if the file exists, open it
		f1, err1 = os.OpenFile(pubFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	} else {
		f1, err1 = os.Create(pubFile) // if the file doesn't exists, create it
	}

	check(err1)
	_, err1 = io.WriteString(f1, pubKey)

	// write private key into the second file
	if checkFileIsExist(priFile) { // if the file exists, open it
		f2, err2 = os.OpenFile(priFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	} else {
		f2, err2 = os.Create(priFile) // if the file doesn't exists, create it
	}

	check(err2)
	_, err2 = io.WriteString(f2, priKey)

	fmt.Println("( p,g,g^a mod p )")
	fmt.Println(pubKey)

}
