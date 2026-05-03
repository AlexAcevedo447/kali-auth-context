package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthConfig struct {
	Secret    string
	Issuer    string
	Audience  string
	ClockSkew time.Duration
}

// NewJWTAuthMiddleware validates Bearer JWT locally using HS256.
// It stores normalized auth data in Fiber locals:
// - auth.sub
// - auth.tid
// - auth.email
// - auth.claims
func NewJWTAuthMiddleware(cfg JWTAuthConfig) fiber.Handler {
	clockSkew := cfg.ClockSkew
	if clockSkew <= 0 {
		clockSkew = 30 * time.Second
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing Authorization header"})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid Authorization header format"})
		}

		tokenString := strings.TrimSpace(parts[1])

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		},
			jwt.WithLeeway(clockSkew),
		)
		if err != nil || token == nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid access token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		if cfg.Issuer != "" && !verifyIssuer(claims, cfg.Issuer) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token issuer"})
		}

		if cfg.Audience != "" && !verifyAudience(claims, cfg.Audience) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token audience"})
		}

		sub, _ := claims["sub"].(string)
		tid, _ := claims["tid"].(string)
		email, _ := claims["email"].(string)
		if strings.TrimSpace(sub) == "" || strings.TrimSpace(tid) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token is missing required claims"})
		}

		c.Locals("auth.sub", sub)
		c.Locals("auth.tid", tid)
		c.Locals("auth.email", email)
		c.Locals("auth.claims", claims)

		return c.Next()
	}
}

func verifyAudience(claims jwt.MapClaims, expected string) bool {
	audRaw, ok := claims["aud"]
	if !ok {
		return false
	}

	switch v := audRaw.(type) {
	case string:
		return v == expected
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok && s == expected {
				return true
			}
		}
	}

	return false
}

func verifyIssuer(claims jwt.MapClaims, expected string) bool {
	issRaw, ok := claims["iss"]
	if !ok {
		return false
	}

	iss, ok := issRaw.(string)
	if !ok {
		return false
	}

	return iss == expected
}
