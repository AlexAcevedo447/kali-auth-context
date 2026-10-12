package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenIssuer struct {
	issuer     string
	audience   string
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	ttl        time.Duration
	kid        string
}

func NewAccessTokenIssuer(cfg *config.Config) (*AccessTokenIssuer, error) {
	if cfg.JWTPrivateKeyPath == "" || cfg.JWTPublicKeyPath == "" {
		return nil, errors.New("JWT private/public key path not configured")
	}
	priv, err := loadPrivateKey(cfg.JWTPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	pub, err := loadPublicKey(cfg.JWTPublicKeyPath)
	if err != nil {
		return nil, err
	}
	ttlMinutes := cfg.JWTAccessTokenTTLMinutes
	if ttlMinutes <= 0 {
		ttlMinutes = 60
	}
	kid := cfg.JWTKeyID
	if kid == "" {
		kid = "main-key"
	}
	return &AccessTokenIssuer{
		issuer:     cfg.JWTIssuer,
		audience:   cfg.JWTAudience,
		privateKey: priv,
		publicKey:  pub,
		ttl:        time.Duration(ttlMinutes) * time.Minute,
		kid:        kid,
	}, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid private key PEM")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid public key PEM")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func (i *AccessTokenIssuer) Issue(tenantId identity.TenantId, userId identity.UserId, email string, roles []*identity.Role, permissions []*identity.Permission) (string, int64, error) {
	if i.privateKey == nil {
		return "", 0, errors.New("jwt private key is not configured")
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = i.kid
	signedToken, err := token.SignedString(i.privateKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, int64(i.ttl.Seconds()), nil
}

func (i *AccessTokenIssuer) PublicKey() *rsa.PublicKey {
	return i.publicKey
}

func (i *AccessTokenIssuer) KeyID() string {
	return i.kid
}
