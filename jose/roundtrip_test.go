package jose_test

import (
	"github.com/absurdlab/pkg/jose"
	"github.com/absurdlab/pkg/jose/jwa"
	"github.com/absurdlab/pkg/jose/jwk"
	"github.com/absurdlab/pkg/jose/jwt"
	"github.com/absurdlab/pkg/timeplus"
	"testing"
)

func TestEncoder(t *testing.T) {
	type etc struct {
		Foo string `json:"foo"`
	}

	type Encoder interface {
		CompactSerialize() (string, error)
	}

	type Decoder interface {
		Decode(token string, dest ...interface{}) error
	}

	cases := []struct {
		name   string
		jwks   *jwk.KeySet
		encode func(jwks *jwk.KeySet) Encoder
		decode func(jwks *jwk.KeySet) (Decoder, []interface{})
		assert func(t *testing.T, dest []interface{}, err error)
	}{
		{
			name: "signing only",
			jwks: jwk.NewKeySet(jwk.GenerateSignatureKey("key1", jwa.ES256, 0)),
			encode: func(jwks *jwk.KeySet) Encoder {
				return jose.Encode().
					Claims(
						new(jwt.Claims).
							GenerateID().
							WithAudience("tester").
							WithExpiryInFuture(timeplus.Second(600)).
							WithIssuedAtNow().
							WithSubject("test"),
					).
					Claims(etc{Foo: "bar"}).
					Sign(jwks, jwa.ES256)
			},
			decode: func(jwks *jwk.KeySet) (Decoder, []interface{}) {
				return jose.Decode().Verify(jwks, jwa.ES256), []interface{}{
					new(jwt.Claims),
					new(etc),
				}
			},
			assert: func(t *testing.T, dest []interface{}, err error) {
				if err != nil {
					t.Error(err)
				}
				if got, want := dest[0].(*jwt.Claims).Subject, "test"; got != want {
					t.Errorf("decoded claims mismatch (sub), got %s, want %s", got, want)
				}
				if got, want := dest[1].(*etc).Foo, "bar"; got != want {
					t.Errorf("decoded claims mismatch (foo), got %s, want %s", got, want)
				}
			},
		},
		{
			name: "encryption only",
			jwks: jwk.NewKeySet(jwk.GenerateEncryptionKey("key1", jwa.RSA1_5, 0)),
			encode: func(jwks *jwk.KeySet) Encoder {
				return jose.Encode().
					Claims(
						new(jwt.Claims).
							GenerateID().
							WithAudience("tester").
							WithExpiryInFuture(timeplus.Second(600)).
							WithIssuedAtNow().
							WithSubject("test"),
					).
					Claims(etc{Foo: "bar"}).
					Encrypt(jwks, jwa.RSA1_5, jwa.A128GCM)
			},
			decode: func(jwks *jwk.KeySet) (Decoder, []interface{}) {
				return jose.Decode().Decrypt(jwks, jwa.RSA1_5), []interface{}{
					new(jwt.Claims),
					new(etc),
				}
			},
			assert: func(t *testing.T, dest []interface{}, err error) {
				if err != nil {
					t.Error(err)
				}
				if got, want := dest[0].(*jwt.Claims).Subject, "test"; got != want {
					t.Errorf("decoded claims mismatch (sub), got %s, want %s", got, want)
				}
				if got, want := dest[1].(*etc).Foo, "bar"; got != want {
					t.Errorf("decoded claims mismatch (foo), got %s, want %s", got, want)
				}
			},
		},
		{
			name: "signed and encrypted",
			jwks: jwk.NewKeySet(
				jwk.GenerateSignatureKey("key1", jwa.ES256, 0),
				jwk.GenerateEncryptionKey("key2", jwa.RSA1_5, 0),
			),
			encode: func(jwks *jwk.KeySet) Encoder {
				return jose.Encode().
					Claims(
						new(jwt.Claims).
							GenerateID().
							WithAudience("tester").
							WithExpiryInFuture(timeplus.Second(600)).
							WithIssuedAtNow().
							WithSubject("test"),
					).
					Claims(etc{Foo: "bar"}).
					Sign(jwks, jwa.ES256).
					Encrypt(jwks, jwa.RSA1_5, jwa.A128GCM)
			},
			decode: func(jwks *jwk.KeySet) (Decoder, []interface{}) {
				return jose.Decode().Decrypt(jwks, jwa.RSA1_5).Verify(jwks, jwa.ES256), []interface{}{
					new(jwt.Claims),
					new(etc),
				}
			},
			assert: func(t *testing.T, dest []interface{}, err error) {
				if err != nil {
					t.Error(err)
				}
				if got, want := dest[0].(*jwt.Claims).Subject, "test"; got != want {
					t.Errorf("decoded claims mismatch (sub), got %s, want %s", got, want)
				}
				if got, want := dest[1].(*etc).Foo, "bar"; got != want {
					t.Errorf("decoded claims mismatch (foo), got %s, want %s", got, want)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			token, err := c.encode(c.jwks).CompactSerialize()
			if err != nil {
				t.Error(err)
			}

			t.Log(token)

			decoder, destinations := c.decode(c.jwks)
			err = decoder.Decode(token, destinations...)
			c.assert(t, destinations, err)
		})
	}
}
