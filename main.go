package cryptopuzzle

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/Echowxsy/go-crypto-puzzle/archive"
	"github.com/Echowxsy/go-crypto-puzzle/encryptor"
	"github.com/Echowxsy/go-crypto-puzzle/utils"
)

type Options struct {
	PrimeBits    int
	PrimeRounds  int
	OpsPerSecond int
	Duration     int
}

type CryptoPuzzle struct {
	Options
}

func NewCryptoPuzzle(o Options) *CryptoPuzzle {
	PRIME_BITS := 100
	PRIME_ROUNDS := 6
	OPS_PER_SECOND := 3300000
	DURATION := 1000
	if o.PrimeBits == 0 {
		o.PrimeBits = PRIME_BITS
	}
	if o.PrimeRounds == 0 {
		o.PrimeRounds = PRIME_ROUNDS
	}
	if o.OpsPerSecond == 0 {
		o.OpsPerSecond = OPS_PER_SECOND
	}
	if o.Duration == 0 {
		o.Duration = DURATION
	}
	return &CryptoPuzzle{
		Options: o,
	}
}

func (c *CryptoPuzzle) Generate(message []byte) ([]byte, error) {

	p := utils.RandomPrime(c.PrimeBits, c.PrimeRounds)
	q := utils.RandomPrime(c.PrimeBits, c.PrimeRounds)

	n := new(big.Int).Mul(p, q)
	n1 := new(big.Int).Mul(new(big.Int).Sub(p, big.NewInt(1)), new(big.Int).Sub(q, big.NewInt(1)))

	S := c.OpsPerSecond
	T := c.Duration
	t := new(big.Int).Mul(big.NewInt(int64(S)), big.NewInt(int64(T/1000)))

	K := make([]byte, 32)
	_, err := rand.Read(K)
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	hash.Write(K)
	K = hash.Sum(nil)

	M := message

	Cm, err := encryptor.Encrypt(M, K, nil, 0)

	if err != nil {
		return nil, err
	}

	a, _ := utils.InRange(big.NewInt(1), new(big.Int).Sub(n, big.NewInt(1)))
	e := utils.FastModExp(big.NewInt(2), t, n1)
	b := utils.FastModExp(a, e, n)
	Ck := new(big.Int).Add(new(big.Int).SetBytes(K), b)
	ar := archive.Archive{
		N: n, A: a, T: t, Ck: Ck,
		Cm: Cm}
	arch := archive.ArchiveToByte(ar)

	return arch, nil
}
