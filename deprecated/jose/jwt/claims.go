package jwt

import (
	"github.com/absurdlab/pkg/randstring"
	"github.com/absurdlab/pkg/timeplus"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Claims is an alias of jwt.Claims.
type Claims jwt.Claims

func (c *Claims) GenerateID() *Claims {
	c.ID = randstring.MustHex(16)
	return c
}

func (c *Claims) WithExpiry(exp timeplus.Timestamp) *Claims {
	c.Expiry = jwt.NewNumericDate(exp.Time())
	return c
}

func (c *Claims) WithExpiryInFuture(seconds timeplus.Second) *Claims {
	return c.WithExpiry(timeplus.Now().AddSecond(seconds))
}

func (c *Claims) WithNotBefore(nbf timeplus.Timestamp) *Claims {
	c.NotBefore = jwt.NewNumericDate(nbf.Time())
	return c
}

func (c *Claims) WithIssuedAt(iat timeplus.Timestamp) *Claims {
	c.IssuedAt = jwt.NewNumericDate(iat.Time())
	return c
}

func (c *Claims) WithIssuedAtNow() *Claims {
	return c.WithIssuedAt(timeplus.Now())
}

func (c *Claims) WithIssuer(iss string) *Claims {
	c.Issuer = iss
	return c
}

func (c *Claims) WithAudience(aud ...string) *Claims {
	c.Audience = aud
	return c
}

func (c *Claims) WithSubject(sub string) *Claims {
	c.Subject = sub
	return c
}

// Flatten is a voluntary interface to implement, so that the claim object can be flattened
// and supplied to JWT builder in pieces. This is useful when the custom JWT claim object
// has a hierarchical structure, but wishes to be encoded as a flat JWT claim.
type Flatten interface {
	// Flatten returns a list of objects, which will be supplied to JWT builder in sequence.
	Flatten() []interface{}
}
