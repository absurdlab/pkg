package jwt

import (
	uuid "github.com/satori/go.uuid"
	"gopkg.in/square/go-jose.v2/jwt"
	"time"
)

// Claims is an alias of jwt.Claims.
type Claims jwt.Claims

func (c *Claims) GenerateID() *Claims {
	c.ID = uuid.NewV4().String()
	return c
}

func (c *Claims) WithExpiry(exp time.Time) *Claims {
	c.Expiry = jwt.NewNumericDate(exp)
	return c
}

func (c *Claims) WithExpiryInFuture(seconds int64) *Claims {
	return c.WithExpiry(time.Now().Add(time.Duration(seconds) * time.Second))
}

func (c *Claims) WithNotBefore(nbf time.Time) *Claims {
	c.NotBefore = jwt.NewNumericDate(nbf)
	return c
}

func (c *Claims) WithIssuedAt(iat time.Time) *Claims {
	c.IssuedAt = jwt.NewNumericDate(iat)
	return c
}

func (c *Claims) WithIssuedAtNow() *Claims {
	return c.WithIssuedAt(time.Now())
}

func (c *Claims) WithIssuer(iss string) *Claims {
	c.Issuer = iss
	return c
}

func (c *Claims) WithAudience(aud ...string) *Claims {
	if c.Audience == nil {
		c.Audience = []string{}
	}
	if len(aud) > 0 {
		c.Audience = append(c.Audience, aud...)
	}
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
