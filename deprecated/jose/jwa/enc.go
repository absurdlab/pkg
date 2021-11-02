package jwa

import "gopkg.in/square/go-jose.v2"

const (
	ED25519            = string(jose.ED25519)
	RSA1_5             = string(jose.RSA1_5)
	RSA_OAEP           = string(jose.RSA_OAEP)
	RSA_OAEP_256       = string(jose.RSA_OAEP_256)
	A128KW             = string(jose.A128KW)
	A192KW             = string(jose.A192KW)
	A256KW             = string(jose.A256KW)
	DIRECT             = string(jose.DIRECT)
	ECDH_ES            = string(jose.ECDH_ES)
	ECDH_ES_A128KW     = string(jose.ECDH_ES_A128KW)
	ECDH_ES_A192KW     = string(jose.ECDH_ES_A192KW)
	ECDH_ES_A256KW     = string(jose.ECDH_ES_A256KW)
	A128GCMKW          = string(jose.A128GCMKW)
	A192GCMKW          = string(jose.A192GCMKW)
	A256GCMKW          = string(jose.A256GCMKW)
	PBES2_HS256_A128KW = string(jose.PBES2_HS256_A128KW)
	PBES2_HS384_A192KW = string(jose.PBES2_HS384_A192KW)
	PBES2_HS512_A256KW = string(jose.PBES2_HS512_A256KW)
)

const (
	A128CBC_HS256 = string(jose.A128CBC_HS256)
	A192CBC_HS384 = string(jose.A192CBC_HS384)
	A256CBC_HS512 = string(jose.A256CBC_HS512)
	A128GCM       = string(jose.A128GCM)
	A192GCM       = string(jose.A192GCM)
	A256GCM       = string(jose.A256GCM)
)

var (
	// EnumEncrypt is a collection of all JWA encryption algorithms
	EnumEncrypt = []string{
		ED25519,
		RSA1_5, RSA_OAEP, RSA_OAEP_256,
		A128KW, A192KW, A256KW,
		DIRECT,
		ECDH_ES, ECDH_ES_A128KW, ECDH_ES_A192KW, ECDH_ES_A256KW,
		A128GCMKW, A192GCMKW, A256GCMKW,
		PBES2_HS256_A128KW, PBES2_HS384_A192KW, PBES2_HS512_A256KW,
	}
	// EnumEncode is a collection of all JWA encryption encoding algorithm.
	EnumEncode = []string{
		A128CBC_HS256, A192CBC_HS384, A256CBC_HS512,
		A128GCM, A192GCM, A256GCM,
	}
)
