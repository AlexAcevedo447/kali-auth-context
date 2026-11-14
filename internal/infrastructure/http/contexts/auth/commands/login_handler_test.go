package commands

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	authcommands "kali-auth-context/internal/application/auth/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/config"
	"kali-auth-context/internal/infrastructure/security"
	"kali-auth-context/internal/ports"
)

type plainHasher struct{}

func (plainHasher) Hash(p string) (string, error) { return "hash:" + p, nil }
func (plainHasher) Compare(h, p string) error {
	if h != "hash:"+p {
		return errors.New("mismatch")
	}
	return nil
}
func (plainHasher) NeedsRehash(string) bool { return false }

type recordingCandidates struct {
	all     []ports.LoginCandidate
	gotCode string
}

func (r *recordingCandidates) FindByEmail(_ context.Context, email, code string) ([]ports.LoginCandidate, error) {
	r.gotCode = code
	var out []ports.LoginCandidate
	for _, c := range r.all {
		if c.User.Email == email && (code == "" || c.TenantCode == code) {
			out = append(out, c)
		}
	}
	return out, nil
}

type noRoles struct{}

func (noRoles) GetByUser(context.Context, identity.TenantId, identity.UserId) ([]*identity.UserRole, error) {
	return nil, nil
}
func (noRoles) GetById(context.Context, identity.RoleId) (*identity.Role, error) {
	return nil, errors.New("unused")
}
func (noRoles) GetByRole(context.Context, identity.TenantId, identity.RoleId) ([]*identity.RolePermission, error) {
	return nil, nil
}

type noPerms struct{}

func (noPerms) GetById(context.Context, identity.PermissionId) (*identity.Permission, error) {
	return nil, errors.New("unused")
}

func newTestIssuer(t *testing.T) *security.AccessTokenIssuer {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	priv := filepath.Join(dir, "private.pem")
	pub := filepath.Join(dir, "public.pem")
	if err := os.WriteFile(priv, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pub, pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey)}), 0o600); err != nil {
		t.Fatal(err)
	}
	issuer, err := security.NewAccessTokenIssuer(&config.Config{
		JWTIssuer: "test", JWTAudience: "test", JWTPrivateKeyPath: priv, JWTPublicKeyPath: pub,
	})
	if err != nil {
		t.Fatal(err)
	}
	return issuer
}

func loginApp(t *testing.T, repo *recordingCandidates) *fiber.App {
	cmd := authcommands.NewLoginCommand(repo, noRoles{}, noRoles{}, noRoles{}, noPerms{}, plainHasher{})
	app := fiber.New()
	app.Post("/login", NewLoginHandler(cmd, newTestIssuer(t), platformTenantForTests).Handle)
	return app
}

func cand(tenant identity.TenantId, code string, status identity.TenantStatus) ports.LoginCandidate {
	return ports.LoginCandidate{
		User:         &identity.User{Id: identity.UserId("u-" + code), TenantId: tenant, Email: "ana@example.com", Password: "hash:pw"},
		TenantName:   strings.ToUpper(code),
		TenantCode:   code,
		TenantStatus: status,
	}
}

func postLogin(t *testing.T, app *fiber.App, body string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var out map[string]any
	_ = json.Unmarshal(raw, &out)
	return resp.StatusCode, out
}

func TestLoginHandler_NormalizesTenantCodeAndResolvesTenant(t *testing.T) {
	repo := &recordingCandidates{all: []ports.LoginCandidate{
		cand("tenant-a", "acme", identity.TenantStatusActive),
		cand("tenant-b", "globex", identity.TenantStatusActive),
	}}
	app := loginApp(t, repo)

	status, body := postLogin(t, app, `{"email":"  ana@example.com ","password":"pw","tenant_code":"  GloBex  "}`)
	if status != http.StatusOK {
		t.Fatalf("expected 200, got %d %v", status, body)
	}
	if repo.gotCode != "globex" {
		t.Fatalf("tenant code not normalized (lowercase/trim): %q", repo.gotCode)
	}
	if body["tenant_id"] != "tenant-b" || body["tenant_code"] != "globex" || body["access_token"] == "" {
		t.Fatalf("unexpected response %v", body)
	}
}

func TestLoginHandler_ErrorMapping(t *testing.T) {
	both := []ports.LoginCandidate{
		cand("tenant-a", "acme", identity.TenantStatusActive),
		cand("tenant-b", "globex", identity.TenantStatusActive),
	}
	cases := []struct {
		name       string
		candidates []ports.LoginCandidate
		body       string
		wantStatus int
		wantError  string
	}{
		{"ambiguous without code", both, `{"email":"ana@example.com","password":"pw"}`, http.StatusConflict, "company code required"},
		{"wrong tenant code", both, `{"email":"ana@example.com","password":"pw","tenant_code":"initech"}`, http.StatusUnauthorized, "invalid credentials"},
		{"wrong password", both, `{"email":"ana@example.com","password":"nope"}`, http.StatusUnauthorized, "invalid credentials"},
		{"suspended tenant", []ports.LoginCandidate{cand("tenant-a", "acme", identity.TenantStatusSuspended)}, `{"email":"ana@example.com","password":"pw"}`, http.StatusForbidden, "tenant is suspended"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := loginApp(t, &recordingCandidates{all: tc.candidates})
			status, body := postLogin(t, app, tc.body)
			if status != tc.wantStatus || body["error"] != tc.wantError {
				t.Fatalf("expected %d %q, got %d %v", tc.wantStatus, tc.wantError, status, body)
			}
		})
	}
}

const platformTenantForTests = "tenant-platform"

func TestLoginHandler_IsPlatformTenantFlag(t *testing.T) {
	cases := []struct {
		tenant identity.TenantId
		want   bool
	}{
		{"tenant-platform", true},
		{"tenant-a", false},
	}
	for _, tc := range cases {
		repo := &recordingCandidates{all: []ports.LoginCandidate{cand(tc.tenant, "acme", identity.TenantStatusActive)}}
		status, body := postLogin(t, loginApp(t, repo), `{"email":"ana@example.com","password":"pw"}`)
		if status != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d %v", tc.tenant, status, body)
		}
		if body["is_platform_tenant"] != tc.want {
			t.Fatalf("%s: expected is_platform_tenant=%v, got %v", tc.tenant, tc.want, body["is_platform_tenant"])
		}
	}
}
