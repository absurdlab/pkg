package jose

import (
	"errors"
	"fmt"
	"github.com/absurdlab/pkg/jose/jwa"
	"github.com/absurdlab/pkg/jose/jwk"
	"github.com/absurdlab/pkg/jose/jwt"
	"gopkg.in/square/go-jose.v2"
	squarejwt "gopkg.in/square/go-jose.v2/jwt"
)

var ErrEncode = errors.New("encode token error")

// Encoder abstracts JWS/JWE encoding.
type Encoder interface {
	// CompactSerialize produces a encoded JWT token.
	CompactSerialize() (string, error)
}

// Encode is the entrypoint for configuring the encoding process.
func Encode() *encodeOptions {
	return &encodeOptions{}
}

type encodeOptions struct {
	claims []interface{}

	sigAlg  string
	sigJwks *jwk.KeySet

	encryptAlg  string
	encryptEnc  string
	encryptJwks *jwk.KeySet
}

// Claims appends the given object to the list of claims to be serialized. Multiple calls to Claims appends additional
// claims. If the object implements jwt.Flatten interface, it will be flattened before appending.
func (o *encodeOptions) Claims(c interface{}) *encodeOptions {
	if flat, ok := c.(jwt.Flatten); ok {
		o.claims = append(o.claims, flat.Flatten()...)
	} else {
		o.claims = append(o.claims, c)
	}
	return o
}

// Sign tells the encoding process to sign the claims with the key of given algorithm from the key set.
func (o *encodeOptions) Sign(jwks *jwk.KeySet, alg string) *encodeOptions {
	o.sigJwks = jwks
	o.sigAlg = alg
	return o
}

// Encrypt tells the encoding process to encrypt the claims with the key of given algorithm from the key set.
func (o *encodeOptions) Encrypt(jwks *jwk.KeySet, alg string, enc string) *encodeOptions {
	o.encryptJwks = jwks
	o.encryptAlg = alg
	o.encryptEnc = enc
	return o
}

// CompactSerialize produces the final jwt/jwe token.
func (o *encodeOptions) CompactSerialize() (token string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrEncode, err)
		}
	}()

	err = o.preSerializeCheck()
	if err != nil {
		return
	}

	if o.shouldSign() && !o.shouldEncrypt() {
		return o.buildWithSigningOnly()
	}

	if !o.shouldSign() && o.shouldEncrypt() {
		return o.buildWithEncryptionOnly()
	}

	return o.buildWithSigningAndEncryption()
}

func (o *encodeOptions) preSerializeCheck() error {
	if len(o.claims) == 0 {
		return errors.New("no claims")
	}

	if !o.shouldSign() && !o.shouldEncrypt() {
		return errors.New("no sign nor encrypt required")
	}

	return nil
}

func (o *encodeOptions) signer() (jose.Signer, error) {
	key, err := o.sigJwks.KeyForSigning(o.sigAlg)
	if err != nil {
		return nil, err
	}

	signKey := jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(key.Algorithm),
		Key:       key.Key,
	}

	opts := new(jose.SignerOptions).
		WithHeader("kid", key.KeyID)

	return jose.NewSigner(signKey, opts)
}

func (o *encodeOptions) encrypter(contentIsJWT bool) (jose.Encrypter, error) {
	key, err := o.encryptJwks.KeyForEncryption(o.encryptAlg)
	if err != nil {
		return nil, err
	}

	recipient := jose.Recipient{
		Algorithm: jose.KeyAlgorithm(o.encryptAlg),
		Key:       key.Public().Key,
		KeyID:     key.KeyID,
	}

	opts := new(jose.EncrypterOptions).
		WithHeader("kid", key.KeyID)
	if contentIsJWT {
		opts = opts.WithContentType("JWT")
	}

	return jose.NewEncrypter(jose.ContentEncryption(o.encryptEnc), recipient, opts)
}

func (o *encodeOptions) buildWithSigningOnly() (string, error) {
	signer, err := o.signer()
	if err != nil {
		return "", err
	}

	builder := squarejwt.Signed(signer)
	for _, each := range o.claims {
		builder = builder.Claims(each)
	}

	return builder.CompactSerialize()
}

func (o *encodeOptions) buildWithEncryptionOnly() (string, error) {
	encrypter, err := o.encrypter(false)
	if err != nil {
		return "", err
	}

	builder := squarejwt.Encrypted(encrypter)
	for _, each := range o.claims {
		builder = builder.Claims(each)
	}

	return builder.CompactSerialize()
}

func (o *encodeOptions) buildWithSigningAndEncryption() (string, error) {
	signer, err := o.signer()
	if err != nil {
		return "", err
	}

	encrypter, err := o.encrypter(true)
	if err != nil {
		return "", err
	}

	builder := squarejwt.SignedAndEncrypted(signer, encrypter)
	for _, each := range o.claims {
		builder = builder.Claims(each)
	}

	return builder.CompactSerialize()
}

func (o *encodeOptions) shouldEncrypt() bool {
	return o.encryptJwks != nil && jwa.IsDefined(o.encryptAlg) && jwa.IsDefined(o.encryptEnc)
}

func (o *encodeOptions) shouldSign() bool {
	return o.sigJwks != nil && jwa.IsDefined(o.sigAlg)
}
