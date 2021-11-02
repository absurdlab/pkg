package randstring

import (
	"math/rand"
	"strings"
	"time"
)

// Option describes the (a part) generated string.
type Option struct {
	// Charset is the string containing all candidate characters
	Charset string
	// Count is the number of characters from Charset that will be used
	Count uint
}

var defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate generates a new random string according to the given options. The generated string will have
// a length equalling to the sum of counts in all options, and consists of the exact number of characters
// in each charset.
func Generate(options ...Option) string {
	return GenerateWithRand(defaultRand, options...)
}

// GenerateWithRand is like Generate, but uses a custom rand.Rand.
func GenerateWithRand(rand *rand.Rand, options ...Option) string {
	var sb strings.Builder
	for _, opt := range options {
		for i := 0; i < int(opt.Count); i++ {
			sb.WriteRune(rune(opt.Charset[rand.Intn(len(opt.Charset))]))
		}
	}

	runes := []rune(sb.String())
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	return string(runes)
}
