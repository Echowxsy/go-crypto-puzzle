package cryptopuzzle_test

import (
	"encoding/base64"
	"testing"

	cryptopuzzle "github.com/Echowxsy/go-crypto-puzzle"
)

func TestNewCryptoPuzzle(t *testing.T) {
	options := cryptopuzzle.Options{
		PrimeBits:    100,
		PrimeRounds:  6,
		OpsPerSecond: 3300000,
		Duration:     1000,
	}
	c := cryptopuzzle.NewCryptoPuzzle(options)
	p, err := c.Generate([]byte("Hello, World!"))
	if err != nil {
		t.Error(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(p))
}
