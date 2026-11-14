package middleware

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// --- Driver database/sql falso (solo biblioteca estándar) --------------------
//
// TenantGuards consulta la base de datos con *bun.DB. Para probarlo sin Postgres
// se registra un driver en memoria que responde a las dos consultas que hace el
// middleware: dueño (tenant_id) de un recurso y existencia de un tenant.

// fakeOwners: tabla -> id -> tenant dueño. Para "tenants" basta con la clave.
var (
	fakeOwnersMu sync.Mutex
	fakeOwners   = map[string]map[string]string{}
)

var idPattern = regexp.MustCompile(`"?id"? = '([^']*)'`)

type fakeDriver struct{}

func (fakeDriver) Open(string) (driver.Conn, error) { return fakeConn{}, nil }

type fakeConn struct{}

func (fakeConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("prepare not supported") }
func (fakeConn) Close() error                        { return nil }
func (fakeConn) Begin() (driver.Tx, error)           { return nil, errors.New("tx not supported") }

func (fakeConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	m := idPattern.FindStringSubmatch(query)
	if m == nil {
		return nil, errors.New("fake driver: unexpected query: " + query)
	}
	id := m[1]
	fakeOwnersMu.Lock()
	defer fakeOwnersMu.Unlock()

	if strings.Contains(strings.ToUpper(query), "SELECT EXISTS") {
		_, ok := fakeOwners["tenants"][id]
		return &fakeRows{cols: []string{"exists"}, vals: [][]driver.Value{{ok}}}, nil
	}
	for table, owners := range fakeOwners {
		if strings.Contains(query, table) {
			if owner, ok := owners[id]; ok {
				return &fakeRows{cols: []string{"tenant_id"}, vals: [][]driver.Value{{owner}}}, nil
			}
		}
	}
	return &fakeRows{cols: []string{"tenant_id"}}, nil // sin filas -> sql.ErrNoRows
}

type fakeRows struct {
	cols []string
	vals [][]driver.Value
	i    int
}

func (r *fakeRows) Columns() []string { return r.cols }
func (r *fakeRows) Close() error      { return nil }
func (r *fakeRows) Next(dest []driver.Value) error {
	if r.i >= len(r.vals) {
		return io.EOF
	}
	copy(dest, r.vals[r.i])
	r.i++
	return nil
}

var registerOnce sync.Once

const (
	platformTenant = "tenant-platform"
	tenantA        = "tenant-a"
	tenantB        = "tenant-b"
)

func newTestGuards(t *testing.T) *TenantGuards {
	t.Helper()
	registerOnce.Do(func() { sql.Register("tenantguards-fake", fakeDriver{}) })
	fakeOwnersMu.Lock()
	fakeOwners = map[string]map[string]string{
		"tenants":     {platformTenant: "", tenantA: "", tenantB: ""},
		"roles":       {"role-a": tenantA, "role-b": tenantB},
		"user_models": {"user-a": tenantA, "user-b": tenantB},
	}
	fakeOwnersMu.Unlock()
	sqldb, err := sql.Open("tenantguards-fake", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqldb.Close() })
	return NewTenantGuards(bun.NewDB(sqldb, pgdialect.New()), platformTenant)
}

// --- App de prueba -------------------------------------------------------------

// newGuardedApp simula el JWT ya validado (auth.tid) con la cabecera X-Test-Tid y
// monta rutas equivalentes a las del router real. El handler final devuelve el
// tenant efectivo y el cuerpo recibido para comprobar la reescritura.
func newGuardedApp(g *TenantGuards) *fiber.App {
	app := fiber.New()
	fakeJWT := func(c *fiber.Ctx) error {
		if tid := c.Get("X-Test-Tid"); tid != "" {
			c.Locals("auth.tid", tid)
		}
		return c.Next()
	}
	echo := func(c *fiber.Ctx) error {
		var body map[string]any
		_ = json.Unmarshal(c.Body(), &body)
		return c.JSON(fiber.Map{
			"effective": c.Locals(localsEffectiveTenant),
			"query":     c.Query("tenant_id"),
			"body":      body,
		})
	}
	protected := app.Group("", fakeJWT, g.Scope())
	protected.Get("/users", echo)
	protected.Post("/users", echo)
	protected.Get("/users/:userId", g.Owned(OwnedRef{Source: "param", Field: "userId", Table: "user_models"}), echo)
	protected.Post("/rbac/users/roles/assign", g.Owned(
		OwnedRef{Source: "body", Field: "user_id", Table: "user_models"},
		OwnedRef{Source: "body", Field: "role_id", Table: "roles"},
	), echo)
	protected.Get("/rbac/users/roles", g.Owned(OwnedRef{Source: "query", Field: "user_id", Table: "user_models"}), echo)
	tenants := protected.Group("/tenants", g.PlatformOnly())
	tenants.Get("", echo)
	return app
}

func call(t *testing.T, app *fiber.App, method, target, tid, body string) (int, map[string]any) {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if tid != "" {
		req.Header.Set("X-Test-Tid", tid)
	}
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

// --- Scope ---------------------------------------------------------------------

func TestScope_MissingTenantClaim_Unauthorized(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, _ := call(t, app, http.MethodGet, "/users", "", "")
	if status != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", status)
	}
}

func TestScope_QueryTenantDifferentFromToken_Forbidden(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodGet, "/users?tenant_id="+tenantB, tenantA, "")
	if status != http.StatusForbidden || body["error"] != "tenant mismatch" {
		t.Fatalf("expected 403 tenant mismatch, got %d %v", status, body)
	}
}

func TestScope_BodyTenantDifferentFromToken_Forbidden(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodPost, "/users", tenantA, `{"tenant_id":"`+tenantB+`","email":"x@y.z"}`)
	if status != http.StatusForbidden || body["error"] != "tenant mismatch" {
		t.Fatalf("expected 403 tenant mismatch, got %d %v", status, body)
	}
}

func TestScope_QueryTenantEqualToToken_Passes(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodGet, "/users?tenant_id="+tenantA, tenantA, "")
	if status != http.StatusOK || body["effective"] != tenantA || body["query"] != tenantA {
		t.Fatalf("expected 200 scoped to tenant-a, got %d %v", status, body)
	}
}

func TestScope_BodyTenantEqualToToken_Passes(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodPost, "/users", tenantA, `{"tenant_id":"`+tenantA+`","email":"x@y.z"}`)
	if status != http.StatusOK || body["effective"] != tenantA {
		t.Fatalf("expected 200 scoped to tenant-a, got %d %v", status, body)
	}
}

func TestScope_NoTenantSent_InjectsTokenTenant(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodPost, "/users", tenantA, `{"email":"x@y.z"}`)
	if status != http.StatusOK || body["query"] != tenantA {
		t.Fatalf("expected 200 with injected tenant, got %d %v", status, body)
	}
	if b, _ := body["body"].(map[string]any); b["tenant_id"] != tenantA {
		t.Fatalf("tenant_id not injected into body: %v", body["body"])
	}
}

func TestScope_PlatformCanTargetExistingTenant(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodGet, "/users?tenant_id="+tenantB, platformTenant, "")
	if status != http.StatusOK || body["effective"] != tenantB {
		t.Fatalf("expected 200 scoped to tenant-b, got %d %v", status, body)
	}
}

func TestScope_PlatformTargetingUnknownTenant_NotFound(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodGet, "/users?tenant_id=tenant-ghost", platformTenant, "")
	if status != http.StatusNotFound || body["error"] != "resource not found" {
		t.Fatalf("expected 404, got %d %v", status, body)
	}
}

// --- Owned ---------------------------------------------------------------------

func TestOwned_ResourceOfOtherTenant_NotFound(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	cases := []struct{ name, method, target, body string }{
		{"param", http.MethodGet, "/users/user-b", ""},
		{"query", http.MethodGet, "/rbac/users/roles?user_id=user-b", ""},
		{"body", http.MethodPost, "/rbac/users/roles/assign", `{"user_id":"user-a","role_id":"role-b"}`},
		{"missing resource", http.MethodGet, "/users/user-ghost", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := call(t, app, tc.method, tc.target, tenantA, tc.body)
			if status != http.StatusNotFound || body["error"] != "resource not found" {
				t.Fatalf("expected 404 resource not found, got %d %v", status, body)
			}
		})
	}
}

func TestOwned_ResourceOfSameTenant_Passes(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	if status, body := call(t, app, http.MethodGet, "/users/user-a", tenantA, ""); status != http.StatusOK {
		t.Fatalf("param: expected 200, got %d %v", status, body)
	}
	if status, body := call(t, app, http.MethodPost, "/rbac/users/roles/assign", tenantA, `{"user_id":"user-a","role_id":"role-a"}`); status != http.StatusOK {
		t.Fatalf("body: expected 200, got %d %v", status, body)
	}
}

// --- PlatformOnly --------------------------------------------------------------

func TestPlatformOnly_NonPlatformTenant_Forbidden(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodGet, "/tenants", tenantA, "")
	if status != http.StatusForbidden || body["error"] != "platform administrator required" {
		t.Fatalf("expected 403 platform administrator required, got %d %v", status, body)
	}
}

func TestPlatformOnly_PlatformTenant_Passes(t *testing.T) {
	app := newGuardedApp(newTestGuards(t))
	status, body := call(t, app, http.MethodGet, "/tenants", platformTenant, "")
	if status != http.StatusOK {
		t.Fatalf("expected 200, got %d %v", status, body)
	}
}

func TestPlatformOnly_EmptyPlatformConfig_ForbidsEveryone(t *testing.T) {
	g := newTestGuards(t)
	g.platformTenantID = ""
	app := newGuardedApp(g)
	status, _ := call(t, app, http.MethodGet, "/tenants", platformTenant, "")
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 when platform tenant is not configured, got %d", status)
	}
}
