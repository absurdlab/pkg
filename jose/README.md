# jose

Wrapper around `gopkg.in/square/go-jose.v2` for convenient encoding and decoding JWT/JWE needed in OIDC applications.

## Install

```bash
go get -u github.com/absurdlab/pkg/jose
```

## Usage

```go
// encoding
token, _ := jose.Encode().Claims(
    new(jwt.Claims).
        GenerateID().
        WithAudience("tester").
        WithExpiryInFuture(600).
        WithIssuedAtNow().
        WithSubject("test"),
    ).
    Claims(etc{Foo: "bar"}).
    Sign(jwks, jwa.ES256).
	CompactSerialize()

// decoding
var (
	standardClaims = new(jwt.Claims)
	extraClaims = new(etc)
)
_ = jose.Decode().
	Decrypt(jwks, jwa.RSA1_5).
	Verify(jwks, jwa.ES256).
	Decode(token, standardClaims, extraClaims)
```
