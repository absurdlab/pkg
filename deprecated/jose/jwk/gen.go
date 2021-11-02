package jwk

import (
	"github.com/absurdlab/pkg/jose/jwk/internal"
	"gopkg.in/square/go-jose.v2"
)

// GenerateSignatureKey generates a signature Key with the given kid and algorithm.
func GenerateSignatureKey(kid string, alg string, bits int) *Key {
	_, privKey, err := internal.KeygenSig(jose.SignatureAlgorithm(alg), bits)
	if err != nil {
		panic(err)
	}

	return (*Key)(&jose.JSONWebKey{
		Key:       privKey,
		KeyID:     kid,
		Algorithm: alg,
		Use:       UseSig,
	})
}

// GenerateEncryptionKey generates an encryption Key with the given kid and algorithm.
func GenerateEncryptionKey(kid string, alg string, bits int) *Key {
	_, privKey, err := internal.KeygenEnc(jose.KeyAlgorithm(alg), bits)
	if err != nil {
		panic(err)
	}

	return (*Key)(&jose.JSONWebKey{
		Key:       privKey,
		KeyID:     kid,
		Algorithm: alg,
		Use:       UseEnc,
	})
}
