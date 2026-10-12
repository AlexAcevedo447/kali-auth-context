package security

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"

	"github.com/gofiber/fiber/v2"
)

type jwk struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

func NewJWKSHandler(pub *rsa.PublicKey, kid string) fiber.Handler {
	       return func(c *fiber.Ctx) error {
		       jwkKey := jwk{
			       Kty: "RSA",
			       Alg: "RS256",
			       Use: "sig",
			       Kid: kid,
			       N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
			       E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
		       }
		       jwks := jwks{Keys: []jwk{jwkKey}}
		       c.Set("Content-Type", "application/json")
		       return c.Status(fiber.StatusOK).JSON(jwks)
	       }
}
