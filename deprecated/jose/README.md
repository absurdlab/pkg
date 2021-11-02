# jose

[![.github/workflows/jose.yaml](https://github.com/absurdlab/pkg/actions/workflows/jose.yaml/badge.svg)](https://github.com/absurdlab/pkg/actions/workflows/jose.yaml)

Wrapper library for github.com/square/go-jose.v2 for OpenID Connect cryptography needs.

```shell
go get -u github.com/absurdlab/pkg/jose
```

## Features

- :white_check_mark: Complete list of JWA algorithms
- :key: Convenient JWK and JWKS wrappers
- :pouch: Convenient JWT claims wrapper
- :left_right_arrow: Streamline encoding and decoding

## Usage

```go
// generate a key set
jwks := jwk.NewKeySet(
    jwk.GenerateSignatureKey("key1", jwa.ES256, 0),
    jwk.GenerateEncryptionKey("key2", jwa.RSA1_5, 0),
),

// create a standard claim
claims := new(jwt.Claims).
    GenerateID().
    WithIssuer("https://absurdlab.io").
    WithSubject("test").
    WithAudience("https://test.org").
    WithIssuedAtNow().
    WithExpiryInFuture(timeplus.Second(600))
	
// encode it
token, _ := jose.Encode().
    Claims(claims).
    Claims(someOtherClaims).
    Sign(myJwks, jwa.ES256).
    Encrypt(clientJwks, jwa.RSA1_5, jwa.A128GCM).
    CompactSerialize()

// decode it
claims, someOtherClaims := new(jwt.Claims), new(SomeOtherClaims)
_ = jose.Decode().
    Verify(serverJwks).
    Decrypt(myJwks).
    Decode(token, claims, someOtherClaims)
```