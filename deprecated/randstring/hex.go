package randstring

import (
	"encoding/hex"
	"math/rand"
)

// Hex generates a random hex string of n bytes.
func Hex(n int) (string, error) {
	return HexWithRand(defaultRand, n)
}

// HexWithRand is Hex, but uses a custom rand.Rand
func HexWithRand(r *rand.Rand, n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := r.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// MustHex is Hex, but panics on error
func MustHex(n int) string {
	str, err := Hex(n)
	if err != nil {
		panic(err)
	}
	return str
}

// MustHexWithRand is HexWithRand, but panics on error
func MustHexWithRand(r *rand.Rand, n int) string {
	str, err := HexWithRand(r, n)
	if err != nil {
		panic(err)
	}
	return str
}
