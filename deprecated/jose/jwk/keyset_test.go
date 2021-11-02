package jwk_test

import (
	"github.com/absurdlab/pkg/jose/jwa"
	"github.com/absurdlab/pkg/jose/jwk"
	"os"
	"testing"
)

func TestReadKeySet(t *testing.T) {
	f, err := os.Open("testdata/jwks.json")
	if err != nil {
		t.Error(err)
	}

	jwks, err := jwk.ReadKeySet(f)
	if err != nil {
		t.Error(err)
	}

	if want, got := 2, jwks.Count(); want != got {
		t.Errorf("jwks count error, want %d, got %d", want, got)
	}
}

func TestKeySet_KeyForSigning(t *testing.T) {
	jwks := jwk.NewKeySet(
		jwk.GenerateSignatureKey("1", jwa.ES256, 0),
		jwk.GenerateSignatureKey("2", jwa.RS256, 0),
		jwk.GenerateSignatureKey("3", jwa.ES256, 0),
	)

	key, err := jwks.KeyForSigning(jwa.ES256)
	if err != nil {
		t.Error(err)
	}

	if want, got := "3", key.KeyID; want != got {
		t.Errorf("key id error, want %s, got %s", want, got)
	}
}
