package commands

import (
	"context"
	"errors"
	"testing"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

// --- Fakes de los puertos ---------------------------------------------------

// fakeHasher: el "hash" de una contraseña es "hash:" + contraseña.
type fakeHasher struct{}

func (fakeHasher) Hash(p string) (string, error) { return "hash:" + p, nil }
func (fakeHasher) Compare(h, p string) error {
	if h != "hash:"+p {
		return errors.New("mismatch")
	}
	return nil
}
func (fakeHasher) NeedsRehash(string) bool { return false }

// fakeCandidates emula el filtro SQL real: sin código devuelve todas las cuentas
// del correo; con código, solo las de la empresa con ese código exacto.
type fakeCandidates struct {
	all         []ports.LoginCandidate
	gotEmail    string
	gotCode     string
	calls       int
	errToReturn error
}

func (f *fakeCandidates) FindByEmail(_ context.Context, email, tenantCode string) ([]ports.LoginCandidate, error) {
	f.calls++
	f.gotEmail, f.gotCode = email, tenantCode
	if f.errToReturn != nil {
		return nil, f.errToReturn
	}
	out := []ports.LoginCandidate{}
	for _, c := range f.all {
		if c.User.Email != email {
			continue
		}
		if tenantCode != "" && c.TenantCode != tenantCode {
			continue
		}
		out = append(out, c)
	}
	return out, nil
}

type fakeUserRoles struct {
	byUser map[identity.UserId][]*identity.UserRole
}

func (f fakeUserRoles) GetByUser(_ context.Context, _ identity.TenantId, u identity.UserId) ([]*identity.UserRole, error) {
	return f.byUser[u], nil
}

type fakeRoles struct {
	byId map[identity.RoleId]*identity.Role
}

func (f fakeRoles) GetById(_ context.Context, id identity.RoleId) (*identity.Role, error) {
	r, ok := f.byId[id]
	if !ok {
		return nil, errors.New("role not found")
	}
	return r, nil
}

type fakeRolePerms struct {
	byRole map[identity.RoleId][]*identity.RolePermission
}

func (f fakeRolePerms) GetByRole(_ context.Context, _ identity.TenantId, id identity.RoleId) ([]*identity.RolePermission, error) {
	return f.byRole[id], nil
}

type fakePerms struct {
	byId map[identity.PermissionId]*identity.Permission
}

func (f fakePerms) GetById(_ context.Context, id identity.PermissionId) (*identity.Permission, error) {
	p, ok := f.byId[id]
	if !ok {
		return nil, errors.New("permission not found")
	}
	return p, nil
}

// --- Datos de prueba ---------------------------------------------------------

const (
	email    = "ana@example.com"
	password = "S3cret!pass"
)

func candidate(tenant identity.TenantId, code string, userId identity.UserId, pw string, status identity.TenantStatus) ports.LoginCandidate {
	return ports.LoginCandidate{
		User: &identity.User{
			Id:       userId,
			TenantId: tenant,
			Email:    email,
			Password: "hash:" + pw,
		},
		TenantName:   "Empresa " + code,
		TenantCode:   code,
		TenantStatus: status,
	}
}

func newLogin(repo *fakeCandidates) *LoginCommand {
	roles := fakeRoles{byId: map[identity.RoleId]*identity.Role{
		"role-a":     {Id: "role-a", TenantId: "tenant-a", Name: "admin"},
		"role-b":     {Id: "role-b", TenantId: "tenant-b", Name: "viewer"},
		"role-alien": {Id: "role-alien", TenantId: "tenant-b", Name: "alien"},
	}}
	userRoles := fakeUserRoles{byUser: map[identity.UserId][]*identity.UserRole{
		"user-a": {{TenantId: "tenant-a", UserId: "user-a", RoleId: "role-a"}, {TenantId: "tenant-a", UserId: "user-a", RoleId: "role-alien"}},
		"user-b": {{TenantId: "tenant-b", UserId: "user-b", RoleId: "role-b"}},
	}}
	rolePerms := fakeRolePerms{byRole: map[identity.RoleId][]*identity.RolePermission{
		"role-a": {{TenantId: "tenant-a", RoleId: "role-a", PermissionId: "perm-a"}},
	}}
	perms := fakePerms{byId: map[identity.PermissionId]*identity.Permission{
		"perm-a": {Id: "perm-a", TenantId: "tenant-a", Resource: "users", Action: "read"},
	}}
	return NewLoginCommand(repo, userRoles, roles, rolePerms, perms, fakeHasher{})
}

// --- Pruebas -----------------------------------------------------------------

func TestLogin_SingleCandidate_ResolvesTenant(t *testing.T) {
	repo := &fakeCandidates{all: []ports.LoginCandidate{
		candidate("tenant-a", "acme", "user-a", password, identity.TenantStatusActive),
	}}

	got, err := newLogin(repo).Execute(&LoginDto{Email: email, Password: password})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TenantId != "tenant-a" || got.TenantCode != "acme" || got.TenantName != "Empresa acme" || got.UserId != "user-a" {
		t.Fatalf("tenant not resolved from account: %+v", got)
	}
	// Defensa en profundidad: el rol de otra empresa asignado al usuario se ignora.
	if len(got.Roles) != 1 || got.Roles[0].Id != "role-a" {
		t.Fatalf("expected only same-tenant role role-a, got %+v", got.Roles)
	}
	if len(got.Permissions) != 1 || got.Permissions[0].Id != "perm-a" {
		t.Fatalf("expected permission perm-a, got %+v", got.Permissions)
	}
	if repo.gotCode != "" {
		t.Fatalf("no tenant code should be sent to repo, got %q", repo.gotCode)
	}
}

func TestLogin_SameCredentialsInTwoTenants_WithoutCode_RequiresCompanyCode(t *testing.T) {
	repo := &fakeCandidates{all: []ports.LoginCandidate{
		candidate("tenant-a", "acme", "user-a", password, identity.TenantStatusActive),
		candidate("tenant-b", "globex", "user-b", password, identity.TenantStatusActive),
	}}

	_, err := newLogin(repo).Execute(&LoginDto{Email: email, Password: password})
	if !errors.Is(err, identity.ErrTenantCodeRequired) {
		t.Fatalf("expected ErrTenantCodeRequired, got %v", err)
	}
	if err.Error() != "company code required" {
		t.Fatalf("unexpected message %q", err.Error())
	}
}

func TestLogin_SameEmailTwoTenants_OnlyOnePasswordMatches_ResolvesThatTenant(t *testing.T) {
	repo := &fakeCandidates{all: []ports.LoginCandidate{
		candidate("tenant-a", "acme", "user-a", "other-password", identity.TenantStatusActive),
		candidate("tenant-b", "globex", "user-b", password, identity.TenantStatusActive),
	}}

	got, err := newLogin(repo).Execute(&LoginDto{Email: email, Password: password})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TenantId != "tenant-b" {
		t.Fatalf("expected tenant-b, got %s", got.TenantId)
	}
}

func TestLogin_WithTenantCode_PicksThatTenant(t *testing.T) {
	repo := &fakeCandidates{all: []ports.LoginCandidate{
		candidate("tenant-a", "acme", "user-a", password, identity.TenantStatusActive),
		candidate("tenant-b", "globex", "user-b", password, identity.TenantStatusActive),
	}}

	got, err := newLogin(repo).Execute(&LoginDto{Email: email, Password: password, TenantCode: "globex"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TenantId != "tenant-b" || got.UserId != "user-b" || got.TenantCode != "globex" {
		t.Fatalf("expected tenant-b/user-b, got %+v", got)
	}
	if repo.gotCode != "globex" {
		t.Fatalf("tenant code not forwarded to repository, got %q", repo.gotCode)
	}
}

func TestLogin_WrongTenantCode_InvalidCredentials(t *testing.T) {
	repo := &fakeCandidates{all: []ports.LoginCandidate{
		candidate("tenant-a", "acme", "user-a", password, identity.TenantStatusActive),
		candidate("tenant-b", "globex", "user-b", password, identity.TenantStatusActive),
	}}

	_, err := newLogin(repo).Execute(&LoginDto{Email: email, Password: password, TenantCode: "initech"})
	if !errors.Is(err, identity.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_WrongPassword_InvalidCredentials(t *testing.T) {
	repo := &fakeCandidates{all: []ports.LoginCandidate{
		candidate("tenant-a", "acme", "user-a", password, identity.TenantStatusActive),
		candidate("tenant-b", "globex", "user-b", password, identity.TenantStatusActive),
	}}

	_, err := newLogin(repo).Execute(&LoginDto{Email: email, Password: "wrong"})
	if !errors.Is(err, identity.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UnknownEmail_InvalidCredentials(t *testing.T) {
	repo := &fakeCandidates{}
	_, err := newLogin(repo).Execute(&LoginDto{Email: "nobody@example.com", Password: password})
	if !errors.Is(err, identity.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_MissingFields_InvalidCredentialsWithoutQuery(t *testing.T) {
	for _, dto := range []*LoginDto{nil, {Email: email}, {Password: password}} {
		repo := &fakeCandidates{}
		_, err := newLogin(repo).Execute(dto)
		if !errors.Is(err, identity.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials for %+v, got %v", dto, err)
		}
		if repo.calls != 0 {
			t.Fatalf("repository must not be queried for incomplete credentials")
		}
	}
}

func TestLogin_SuspendedTenant_Forbidden(t *testing.T) {
	repo := &fakeCandidates{all: []ports.LoginCandidate{
		candidate("tenant-a", "acme", "user-a", password, identity.TenantStatusSuspended),
	}}

	_, err := newLogin(repo).Execute(&LoginDto{Email: email, Password: password})
	if !errors.Is(err, identity.ErrTenantSuspended) {
		t.Fatalf("expected ErrTenantSuspended, got %v", err)
	}
}

// Comportamiento actual (documentado, no es un error según la arquitectura):
// si las credenciales coinciden en una empresa activa y en una suspendida, sin
// código se pide primero el código de empresa (409). Con el código de la
// suspendida se responde 403; con el de la activa, el login funciona.
func TestLogin_SuspendedAmongMultipleCandidates_AsksForCodeFirst(t *testing.T) {
	repo := &fakeCandidates{all: []ports.LoginCandidate{
		candidate("tenant-a", "acme", "user-a", password, identity.TenantStatusSuspended),
		candidate("tenant-b", "globex", "user-b", password, identity.TenantStatusActive),
	}}
	cmd := newLogin(repo)

	if _, err := cmd.Execute(&LoginDto{Email: email, Password: password}); !errors.Is(err, identity.ErrTenantCodeRequired) {
		t.Fatalf("without code: expected ErrTenantCodeRequired, got %v", err)
	}
	if _, err := cmd.Execute(&LoginDto{Email: email, Password: password, TenantCode: "acme"}); !errors.Is(err, identity.ErrTenantSuspended) {
		t.Fatalf("suspended code: expected ErrTenantSuspended, got %v", err)
	}
	got, err := cmd.Execute(&LoginDto{Email: email, Password: password, TenantCode: "globex"})
	if err != nil || got.TenantId != "tenant-b" {
		t.Fatalf("active code: expected tenant-b, got %+v / %v", got, err)
	}
}

func TestLogin_RepositoryError_Propagates(t *testing.T) {
	boom := errors.New("db down")
	repo := &fakeCandidates{errToReturn: boom}
	_, err := newLogin(repo).Execute(&LoginDto{Email: email, Password: password})
	if !errors.Is(err, boom) {
		t.Fatalf("expected repository error, got %v", err)
	}
}
