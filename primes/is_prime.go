package primes

import "math/big"

func IsPrime(number int) bool {
	/*
		Since according to docs "ProbablyPrime is 100% accurate for inputs less than 2⁶⁴.",
		it doesn't matter much which n is chosen.
		Choosing n=1, just in case someone thinks it's a good idea to run this in a version of Go older than 1.8.
	*/
	return big.NewInt(int64(number)).ProbablyPrime(1)
}
