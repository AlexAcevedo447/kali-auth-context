package commands

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	tenantcommands "kali-auth-context/internal/application/tenant/commands"
	"kali-auth-context/internal/domain/identity"
)

type fakeUpdateRepo struct {
	err error
	got *identity.Tenant
}

func (f *fakeUpdateRepo) Update(_ context.Context, t *identity.Tenant) error {
	f.got = t
	return f.err
}

func doUpdate(t *testing.T, repo *fakeUpdateRepo, body string) (int, map[string]any) {
	t.Helper()
	app := fiber.New()
	app.Put("/tenants/:tenantId", NewUpdateHandler(tenantcommands.NewUpdateTenantCommand(repo)).Handle)
	req := httptest.NewRequest(http.MethodPut, "/tenants/tenant-b", strings.NewReader(body))
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

func TestUpdateTenantHandler_DuplicateCode_Returns409(t *testing.T) {
	repo := &fakeUpdateRepo{err: identity.ErrTenantCodeTaken}
	status, body := doUpdate(t, repo, `{"name":"Globex","code":"acme"}`)
	if status != http.StatusConflict || body["error"] != "tenant code already exists" {
		t.Fatalf("expected 409 tenant code already exists, got %d %v", status, body)
	}
	if repo.got == nil || repo.got.Code != "acme" {
		t.Fatalf("code not forwarded to repository: %+v", repo.got)
	}
}

func TestUpdateTenantHandler_OK(t *testing.T) {
	repo := &fakeUpdateRepo{}
	status, _ := doUpdate(t, repo, `{"name":"Globex","code":"globex"}`)
	if status != http.StatusOK {
		t.Fatalf("expected 200, got %d", status)
	}
}

func TestUpdateTenantHandler_InvalidCode_Returns400(t *testing.T) {
	repo := &fakeUpdateRepo{}
	status, _ := doUpdate(t, repo, `{"name":"Globex","code":"Not Valid!"}`)
	if status != http.StatusBadRequest || repo.got != nil {
		t.Fatalf("expected 400 without reaching repository, got %d", status)
	}
}
