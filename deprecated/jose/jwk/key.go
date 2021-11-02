package jwk

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"gopkg.in/square/go-jose.v2"
)

// Standard key use parameters.
const (
	UseSig = "sig"
	UseEnc = "enc"
)

// Key is an alias to jose.JSONWebKey to add more functions.
type Key jose.JSONWebKey

// IsSymmetric returns true if the underlying key uses symmetric algorithms (i.e. HS256)
func (k *Key) IsSymmetric() bool {
	_, ok := k.Key.([]byte)
	return ok
}

// IsPublicAsymmetric returns true if the underlying asymmetric key only has the public portion. If called is not
// sure about the key's symmetry, check with IsSymmetric first.
func (k *Key) IsPublicAsymmetric() bool {
	switch k.Key.(type) {
	case ed25519.PublicKey, *ecdsa.PublicKey, *rsa.PublicKey:
		return true
	default:
		return false
	}
}

// Public returns a new Key with only the public portion of the underlying key. This method
// shall return the same Key if the key is already public or is symmetric.
func (k *Key) Public() *Key {
	if k.IsSymmetric() || k.IsPublicAsymmetric() {
		return k
	}

	return (*Key)(&jose.JSONWebKey{
		Key: func() interface{} {
			raw := k.Key
			switch raw.(type) {
			case ed25519.PrivateKey:
				return raw.(ed25519.PrivateKey).Public()
			case *ecdsa.PrivateKey:
				return raw.(*ecdsa.PrivateKey).Public()
			case *rsa.PrivateKey:
				return raw.(*rsa.PrivateKey).Public()
			default:
				panic("public key conversion is not supported for this key type")
			}
		}(),
		KeyID:     k.KeyID,
		Algorithm: k.Algorithm,
		Use:       k.Use,
	})
}
