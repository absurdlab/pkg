package randstring_test

import (
	"github.com/absurdlab/pkg/randstring"
	"testing"
)

func TestHex(t *testing.T) {
	str, err := randstring.Hex(8)
	if err != nil {
		t.Error(err)
	}
	t.Log(str)
}
