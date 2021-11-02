package randstring_test

import (
	"github.com/absurdlab/pkg/randstring"
	"regexp"
	"testing"
)

func TestGenerate(t *testing.T) {
	type VerifiableOption struct {
		randstring.Option
		Capture *regexp.Regexp
	}

	cases := []struct {
		name    string
		options []VerifiableOption
	}{
		{
			name: "lower alpha",
			options: []VerifiableOption{
				{
					Option:  randstring.Charsets.LowerCaseAlpha(8),
					Capture: regexp.MustCompile(`[a-z]`),
				},
			},
		},
		{
			name: "upper alpha",
			options: []VerifiableOption{
				{
					Option:  randstring.Charsets.UpperCaseAlpha(8),
					Capture: regexp.MustCompile(`[A-Z]`),
				},
			},
		},
		{
			name: "numeric",
			options: []VerifiableOption{
				{
					Option:  randstring.Charsets.Numeric(8),
					Capture: regexp.MustCompile(`[0-9]`),
				},
			},
		},
		{
			name: "alphanumeric",
			options: []VerifiableOption{
				{
					Option:  randstring.Charsets.AlphaNumeric(8),
					Capture: regexp.MustCompile(`[a-zA-Z0-9]`),
				},
			},
		},
		{
			name: "non zero numeric",
			options: []VerifiableOption{
				{
					Option:  randstring.Charsets.NonZeroNumeric(8),
					Capture: regexp.MustCompile(`[1-9]`),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var options []randstring.Option
			for _, each := range c.options {
				options = append(options, each.Option)
			}

			str := randstring.Generate(options...)

			for _, each := range c.options {
				var count uint
				for _, char := range str {
					if each.Capture.MatchString(string(char)) {
						count++
					}
				}
				if count != each.Count {
					t.Fail()
				}
			}
		})
	}
}
