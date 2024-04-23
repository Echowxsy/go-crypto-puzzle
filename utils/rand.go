package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func RandBytes(bytes int) []byte {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}

// func RandomPrime() *big.Int {
// 	for {
// 		p := RandBytes(32)
// 		if p != nil {
// 			p[0] |= 0x80
// 			p[len(p)-1] |= 0x01
// 			prime := new(big.Int).SetBytes(p)
// 			if prime.ProbablyPrime(20) {
// 				return prime
// 			}
// 		}
// 	}
// }

func RandomPrime(bits int, k int) *big.Int {
	if k == 0 {
		k = 8
	}

	for {
		p := make([]byte, bits/8)
		_, err := rand.Read(p)
		if err != nil {
			panic(err)
		}

		n := new(big.Int).SetBytes(p)
		if n.Bit(0) == 0 {
			n.Add(n, big.NewInt(1))
		}

		for {
			if n.ProbablyPrime(k) {
				return n
			}
			n.Add(n, big.NewInt(2))
		}
	}
}

func InRange(min, max *big.Int) (*big.Int, error) {

	if min.Cmp(big.NewInt(0)) < 0 || max.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("Negative ranges are not supported")
	}
	if max.Cmp(min) <= 0 {
		return nil, fmt.Errorf("\"max\" must be at least equal to \"min\" plus 1")
	}

	interval := new(big.Int).Sub(max, min)
	interval.Sub(interval, big.NewInt(1))
	intervalBits := interval.BitLen()

	for {
		nr, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(intervalBits)))
		if err != nil {
			return nil, err
		}
		if nr.Cmp(interval) <= 0 {
			result := new(big.Int).Add(nr, min)
			return result, nil
		}
	}
}
