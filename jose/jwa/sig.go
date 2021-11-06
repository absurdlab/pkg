package jwa

import "gopkg.in/square/go-jose.v2"

const (
	HS256 = string(jose.HS256)
	HS384 = string(jose.HS384)
	HS512 = string(jose.HS512)
	RS256 = string(jose.RS256)
	RS384 = string(jose.RS384)
	RS512 = string(jose.RS512)
	PS256 = string(jose.PS256)
	PS384 = string(jose.PS384)
	PS512 = string(jose.PS512)
	ES256 = string(jose.ES256)
	ES384 = string(jose.ES384)
	ES512 = string(jose.ES512)
)

// EnumSig is a collection of all JWA signature algorithm values.
var EnumSig = []string{
	HS256, HS384, HS512,
	RS256, RS384, RS512,
	PS256, PS384, PS512,
	ES256, ES384, ES512,
}
