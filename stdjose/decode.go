package stdjose

import (
	"errors"
	"fmt"
	"github.com/absurdlab/pkg/stdjose/jwk"
	"gopkg.in/square/go-jose.v2"
	squarejwt "gopkg.in/square/go-jose.v2/jwt"
)

var ErrDecode = errors.New("decode token error")

// Decoder abstracts JWS/JWE decoding.
type Decoder interface {
	// Decode decrypts and/or verifies token and unmarshal the object into the list of destinations.
	Decode(token string, dest ...interface{}) error
}

// Decode is the entrypoint to configure the decoding process.
func Decode() *decodeOptions {
	return &decodeOptions{}
}

type decodeOptions struct {
	verifyAlg  string
	verifyJwks *jwk.KeySet

	decryptAlg  string
	decryptJwks *jwk.KeySet
}

// Verify provides verification key set and algorithm. The decoder will perform signature verification.
func (d *decodeOptions) Verify(jwks *jwk.KeySet, alg string) *decodeOptions {
	d.verifyAlg = alg
	d.verifyJwks = jwks
	return d
}

// Decrypt provides decryption key set and algorithm. The decoder will perform token decryption.
func (d *decodeOptions) Decrypt(jwks *jwk.KeySet, alg string) *decodeOptions {
	d.decryptAlg = alg
	d.decryptJwks = jwks
	return d
}

// Decode decodes the token into the list of destinations.
func (d *decodeOptions) Decode(token string, dest ...interface{}) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrDecode, err)
		}
	}()

	decrypt, verify := d.shouldDecrypt(), d.shouldVerify()
	if !decrypt && !verify {
		return errors.New("no decryption nor verification requested")
	}

	switch {
	case decrypt && !verify:
		return d.decryptOnly(token, dest...)
	case !decrypt && verify:
		return d.verifyOnly(token, dest...)
	case decrypt && verify:
		return d.decryptAndVerify(token, dest...)
	default:
		return errors.New("no decryption nor verification requested")
	}
}

func (d *decodeOptions) decryptOnly(token string, dest ...interface{}) error {
	jwt, err := squarejwt.ParseEncrypted(token)
	if err != nil {
		return err
	}

	jsonWebKey, err := d.decryptKey(jwt.Headers)
	if err != nil {
		return err
	}

	return jwt.Claims(jsonWebKey.Key, dest...)
}

func (d *decodeOptions) verifyOnly(token string, dest ...interface{}) error {
	jwt, err := squarejwt.ParseSigned(token)
	if err != nil {
		return err
	}

	jsonWebKey, err := d.verifyKey(jwt.Headers)
	if err != nil {
		return err
	}

	return jwt.Claims(jsonWebKey.Public().Key, dest...)
}

func (d *decodeOptions) decryptAndVerify(token string, dest ...interface{}) error {
	nestedJWT, err := squarejwt.ParseSignedAndEncrypted(token)
	if err != nil {
		return err
	}

	decryptJsonWebKey, err := d.decryptKey(nestedJWT.Headers)
	if err != nil {
		return err
	}

	jwt, err := nestedJWT.Decrypt(decryptJsonWebKey.Key)
	if err != nil {
		return err
	}

	verifyJsonWebKey, err := d.verifyKey(jwt.Headers)
	if err != nil {
		return err
	}

	return jwt.Claims(verifyJsonWebKey.Public().Key, dest...)
}

func (d *decodeOptions) verifyKey(headers []jose.Header) (*jwk.Key, error) {
	kid, alg := d.getKid(headers), d.getAlg(headers)
	if len(kid) == 0 {
		return nil, errors.New("jws missing kid header")
	} else if len(alg) == 0 {
		return nil, errors.New("jws missing alg header")
	}

	jsonWebKey, err := d.verifyJwks.KeyById(kid)
	if err != nil {
		return nil, err
	}

	if alg != jsonWebKey.Algorithm {
		return nil, errors.New("jws token alg header mismatch with key algorithm")
	} else if len(d.verifyAlg) > 0 && alg != d.verifyAlg {
		return nil, fmt.Errorf("verify algorithm mismatch, got '%s', want '%s'", alg, d.verifyAlg)
	}

	return jsonWebKey, nil
}

func (d *decodeOptions) decryptKey(headers []jose.Header) (*jwk.Key, error) {
	kid, alg := d.getKid(headers), d.getAlg(headers)
	if len(kid) == 0 {
		return nil, errors.New("jwe missing kid header")
	} else if len(alg) == 0 {
		return nil, errors.New("jwe missing alg header")
	}

	jsonWebKey, err := d.decryptJwks.KeyById(kid)
	if err != nil {
		return nil, err
	}

	if alg != jsonWebKey.Algorithm {
		return nil, errors.New("jwe token alg header mismatch with key algorithm")
	} else if len(d.decryptAlg) > 0 && alg != d.decryptAlg {
		return nil, fmt.Errorf("decryption algorithm mismatch, got '%s', want '%s'", alg, d.decryptAlg)
	}

	return jsonWebKey, nil
}

func (d *decodeOptions) shouldDecrypt() bool {
	return d.decryptJwks != nil
}

func (d *decodeOptions) shouldVerify() bool {
	return d.verifyJwks != nil
}

func (d *decodeOptions) getKid(headers []jose.Header) string {
	for _, header := range headers {
		if kid := header.KeyID; len(kid) > 0 {
			return kid
		}
	}
	return ""
}

func (d *decodeOptions) getAlg(headers []jose.Header) string {
	for _, header := range headers {
		if alg := header.Algorithm; len(alg) > 0 {
			return alg
		}
	}
	return ""
}
