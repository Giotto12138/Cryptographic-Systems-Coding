package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

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

// calculate a random generator
func generator(p *big.Int, q *big.Int) {

}

func main() {

	var q *big.Int
	var one *big.Int = big.NewInt(1)
	var two *big.Int = big.NewInt(2)

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
	temp := q.Mul(q, two)
	p := temp.Add(temp, one)

	fmt.Println("q: ", q)
	fmt.Println("p: ", p)

}
