package security

import (
	"errors"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenIssuer struct {
	issuer   string
	audience string
	secret   []byte
	ttl      time.Duration
}

func NewAccessTokenIssuer(cfg *config.Config) *AccessTokenIssuer {
	ttlMinutes := cfg.JWTAccessTokenTTLMinutes
	if ttlMinutes <= 0 {
		ttlMinutes = 60
	}

	return &AccessTokenIssuer{
		issuer:   cfg.JWTIssuer,
		audience: cfg.JWTAudience,
		secret:   []byte(cfg.JWTSecret),
		ttl:      time.Duration(ttlMinutes) * time.Minute,
	}
}

func (i *AccessTokenIssuer) Issue(tenantId identity.TenantId, userId identity.UserId, email string, roles []*identity.Role, permissions []*identity.Permission) (string, int64, error) {
	if len(i.secret) == 0 {
		return "", 0, errors.New("jwt secret is not configured")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(i.ttl)

	roleIDs := make([]string, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, string(r.Id))
	}

	permissionClaims := make([]map[string]string, 0, len(permissions))
	for _, p := range permissions {
		permissionClaims = append(permissionClaims, map[string]string{
			"id":       string(p.Id),
			"resource": p.Resource,
			"action":   p.Action,
		})
	}

	claims := jwt.MapClaims{
		"sub":         userId,
		"tid":         tenantId,
		"email":       email,
		"roles":       roleIDs,
		"permissions": permissionClaims,
		"iat":         now.Unix(),
		"nbf":         now.Unix(),
		"exp":         expiresAt.Unix(),
	}

	if i.issuer != "" {
		claims["iss"] = i.issuer
	}
	if i.audience != "" {
		claims["aud"] = i.audience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(i.secret)
	if err != nil {
		return "", 0, err
	}

	return signedToken, int64(i.ttl.Seconds()), nil
}
