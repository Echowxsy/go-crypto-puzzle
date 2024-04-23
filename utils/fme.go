package utils

import (
	"math/big"
)

func FastModExp(a, b, n *big.Int) *big.Int {
	// copy a b n to avoid modifying the original values
	a = new(big.Int).Set(a)
	b = new(big.Int).Set(b)
	n = new(big.Int).Set(n)

	r := big.NewInt(1)
	x := new(big.Int).Mod(a, n)

	for b.Cmp(big.NewInt(0)) > 0 {
		if new(big.Int).And(b, big.NewInt(1)).Cmp(big.NewInt(1)) == 0 {
			r.Mul(r, x)
			r.Mod(r, n)
		}
		x.Mul(x, x)
		x.Mod(x, n)
		b.Rsh(b, 1)
	}
	return r
}
