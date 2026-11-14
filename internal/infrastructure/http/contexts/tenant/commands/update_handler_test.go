package commands

import (
	"context"
	"database/sql"
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

func TestUpdateTenantHandler_InvalidStatus_Returns400(t *testing.T) {
	repo := &fakeUpdateRepo{}
	status, body := doUpdate(t, repo, `{"name":"Globex","status":"DELETED"}`)
	if status != http.StatusBadRequest || body["error"] != "tenant status must be ACTIVE or SUSPENDED" || repo.got != nil {
		t.Fatalf("expected 400 without reaching repository, got %d %v", status, body)
	}
}

func TestUpdateTenantHandler_KeepsSuspendedStatus(t *testing.T) {
	repo := &fakeUpdateRepo{}
	status, _ := doUpdate(t, repo, `{"name":"Globex","status":"SUSPENDED"}`)
	if status != http.StatusOK || repo.got == nil || repo.got.Status != identity.TenantStatusSuspended {
		t.Fatalf("expected 200 with SUSPENDED forwarded, got %d %+v", status, repo.got)
	}
}

func TestUpdateTenantHandler_UnknownTenant_Returns404(t *testing.T) {
	repo := &fakeUpdateRepo{err: sql.ErrNoRows}
	status, body := doUpdate(t, repo, `{"name":"Globex","status":"ACTIVE"}`)
	if status != http.StatusNotFound || body["error"] != "resource not found" {
		t.Fatalf("expected 404 resource not found, got %d %v", status, body)
	}
}

type fakeStatusRepo struct{ err error }

func (f fakeStatusRepo) Activate(context.Context, identity.TenantId) error { return f.err }
func (f fakeStatusRepo) Suspend(context.Context, identity.TenantId) error  { return f.err }

func doStatus(t *testing.T, action string, repo fakeStatusRepo) int {
	t.Helper()
	app := fiber.New()
	app.Post("/tenants/:tenantId/activate", NewActivateHandler(tenantcommands.NewActivateTenantCommand(repo)).Handle)
	app.Post("/tenants/:tenantId/suspend", NewSuspendHandler(tenantcommands.NewSuspendTenantCommand(repo)).Handle)
	resp, err := app.Test(httptest.NewRequest(http.MethodPost, "/tenants/tenant-b/"+action, nil), -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func TestActivateSuspendHandlers_StatusCodes(t *testing.T) {
	for _, action := range []string{"activate", "suspend"} {
		if got := doStatus(t, action, fakeStatusRepo{}); got != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d", action, got)
		}
		// La empresa no existe: el repositorio devuelve sql.ErrNoRows (sin filas afectadas).
		if got := doStatus(t, action, fakeStatusRepo{err: sql.ErrNoRows}); got != http.StatusNotFound {
			t.Fatalf("%s unknown tenant: expected 404, got %d", action, got)
		}
	}
}

func TestUpdateTenantHandler_BlankName_Returns400(t *testing.T) {
	repo := &fakeUpdateRepo{}
	status, body := doUpdate(t, repo, `{"name":"   ","status":"ACTIVE"}`)
	if status != http.StatusBadRequest || body["error"] != "tenant name is required" || repo.got != nil {
		t.Fatalf("expected 400 tenant name is required, got %d %v", status, body)
	}
}

func TestUpdateTenantHandler_TrimsName(t *testing.T) {
	repo := &fakeUpdateRepo{}
	status, _ := doUpdate(t, repo, `{"name":"  Globex  ","status":"ACTIVE"}`)
	if status != http.StatusOK || repo.got == nil || repo.got.Name != "Globex" {
		t.Fatalf("expected trimmed name, got %d %+v", status, repo.got)
	}
}

// La empresa de plataforma no se puede suspender, ni por /suspend ni por PUT.
func TestPlatformTenant_CannotBeSuspended(t *testing.T) {
	const platform = "tenant-b" // la ruta de prueba usa tenant-b
	app := fiber.New()
	updRepo := &fakeUpdateRepo{}
	app.Put("/tenants/:tenantId", NewUpdateHandler(
		tenantcommands.NewUpdateTenantCommand(updRepo).WithPlatformTenant(platform)).Handle)
	app.Post("/tenants/:tenantId/suspend", NewSuspendHandler(
		tenantcommands.NewSuspendTenantCommand(fakeStatusRepo{}).WithPlatformTenant(platform)).Handle)

	resp, err := app.Test(httptest.NewRequest(http.MethodPost, "/tenants/tenant-b/suspend", nil), -1)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusConflict || !strings.Contains(string(raw), "platform tenant cannot be suspended") {
		t.Fatalf("suspend: expected 409, got %d %s", resp.StatusCode, raw)
	}

	req := httptest.NewRequest(http.MethodPut, "/tenants/tenant-b", strings.NewReader(`{"name":"Plataforma","status":"SUSPENDED"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusConflict || updRepo.got != nil {
		t.Fatalf("PUT SUSPENDED: expected 409 without reaching repository, got %d", resp.StatusCode)
	}

	// Editar la plataforma sin suspenderla sigue permitido.
	req = httptest.NewRequest(http.MethodPut, "/tenants/tenant-b", strings.NewReader(`{"name":"Plataforma","status":"ACTIVE"}`))
	req.Header.Set("Content-Type", "application/json")
	if resp, _ = app.Test(req, -1); resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT ACTIVE: expected 200, got %d", resp.StatusCode)
	}
}
