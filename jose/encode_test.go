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
	cases := []struct {
		name    string
		jwks    *jwk.KeySet
		encoder func(jwks *jwk.KeySet) interface{ CompactSerialize() (string, error) }
	}{
		{
			name: "signing only",
			jwks: jwk.NewKeySet(jwk.GenerateSignatureKey("key1", jwa.ES256, 0)),
			encoder: func(jwks *jwk.KeySet) interface{ CompactSerialize() (string, error) } {
				return jose.Encoder().
					Claims(
						new(jwt.Claims).
							GenerateID().
							WithAudience("tester").
							WithExpiryInFuture(timeplus.Second(600)).
							WithIssuedAtNow().
							WithSubject("test"),
					).
					Claims(map[string]interface{}{"foo": "bar"}).
					Sign(jwks, jwa.ES256)
			},
		},
		{
			name: "encryption only",
			jwks: jwk.NewKeySet(jwk.GenerateEncryptionKey("key1", jwa.RSA1_5, 0)),
			encoder: func(jwks *jwk.KeySet) interface{ CompactSerialize() (string, error) } {
				return jose.Encoder().
					Claims(
						new(jwt.Claims).
							GenerateID().
							WithAudience("tester").
							WithExpiryInFuture(timeplus.Second(600)).
							WithIssuedAtNow().
							WithSubject("test"),
					).
					Claims(map[string]interface{}{"foo": "bar"}).
					Encrypt(jwks, jwa.RSA1_5, jwa.A128GCM)
			},
		},
		{
			name: "signed and encrypted",
			jwks: jwk.NewKeySet(
				jwk.GenerateSignatureKey("key1", jwa.ES256, 0),
				jwk.GenerateEncryptionKey("key1", jwa.RSA1_5, 0),
			),
			encoder: func(jwks *jwk.KeySet) interface{ CompactSerialize() (string, error) } {
				return jose.Encoder().
					Claims(
						new(jwt.Claims).
							GenerateID().
							WithAudience("tester").
							WithExpiryInFuture(timeplus.Second(600)).
							WithIssuedAtNow().
							WithSubject("test"),
					).
					Claims(map[string]interface{}{"foo": "bar"}).
					Sign(jwks, jwa.ES256).
					Encrypt(jwks, jwa.RSA1_5, jwa.A128GCM)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			token, err := c.encoder(c.jwks).CompactSerialize()
			if err != nil {
				t.Error(err)
			}
			t.Log(token)
		})
	}
}
