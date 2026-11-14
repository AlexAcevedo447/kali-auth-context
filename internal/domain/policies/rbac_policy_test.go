package policies

import (
	"errors"
	"testing"

	"kali-auth-context/internal/domain/identity"
)

func TestValidatePermission_Catalog(t *testing.T) {
	cases := []struct {
		resource, action string
		wantR, wantA     string
		wantErr          error
	}{
		{" Invoices ", "READ", "invoices", "read", nil},
		{"invoice-items", "manage", "invoice-items", "manage", nil},
		{"rbac", "manage", "rbac", "manage", nil},
		{"tenants", "manage", "tenants", "manage", nil},
		{"invoices", "*", "invoices", "*", nil},
		{"*", "*", "*", "*", nil},
		{"*", "read", "*", "read", nil},
		{"*", "delete", "", "", identity.ErrUnknownAction},
		{"rbac", "read", "", "", identity.ErrUnknownAction},
		{"tenants", "read", "", "", identity.ErrUnknownAction},
		{"payments", "read", "", "", identity.ErrUnknownResource},
		{"  ", "read", "", "", identity.ErrPermissionResource},
		{"roles", " ", "", "", identity.ErrPermissionAction},
	}
	for _, tc := range cases {
		r, a, err := ValidatePermission(tc.resource, tc.action)
		if !errors.Is(err, tc.wantErr) || r != tc.wantR || a != tc.wantA {
			t.Errorf("ValidatePermission(%q,%q) = (%q,%q,%v), want (%q,%q,%v)", tc.resource, tc.action, r, a, err, tc.wantR, tc.wantA, tc.wantErr)
		}
	}
}

func TestCovers(t *testing.T) {
	g := func(r, a string) Grant { return Grant{Resource: r, Action: a} }
	cases := []struct {
		name string
		have []Grant
		want Grant
		ok   bool
	}{
		{"exacto", []Grant{g("roles", "manage")}, g("roles", "manage"), true},
		{"manage incluye read", []Grant{g("roles", "manage")}, g("roles", "read"), true},
		{"read no incluye manage", []Grant{g("roles", "read")}, g("roles", "manage"), false},
		{"comodín de acción", []Grant{g("invoices", "*")}, g("invoices", "manage"), true},
		{"comodín total", []Grant{g("*", "*")}, g("tenants", "manage"), true},
		{"*/* solo lo cubre */*", []Grant{g("invoices", "*"), g("roles", "manage")}, g("*", "*"), false},
		{"*/manage no cubre */*", []Grant{g("*", "manage")}, g("*", "*"), false},
		{"recurso comodín pedido", []Grant{g("*", "manage")}, g("*", "manage"), true},
		{"otro recurso", []Grant{g("roles", "manage")}, g("permissions", "manage"), false},
		{"mayúsculas y espacios", []Grant{g(" ROLES ", "Manage")}, g("roles", "read"), true},
		{"sin permisos", nil, g("roles", "read"), false},
	}
	for _, tc := range cases {
		if got := Covers(tc.have, tc.want); got != tc.ok {
			t.Errorf("%s: Covers = %v, want %v", tc.name, got, tc.ok)
		}
	}
}

func TestCoversAll_NoEscalation(t *testing.T) {
	rbacOnly := []Grant{{"rbac", "manage"}}
	if CoversAll(rbacOnly, []Grant{{"*", "*"}}) {
		t.Fatal("rbac/manage must not be able to grant */*")
	}
	if CoversAll(rbacOnly, []Grant{{"tenants", "manage"}}) {
		t.Fatal("rbac/manage must not be able to grant tenants/manage")
	}
	if !CoversAll([]Grant{{"*", "*"}}, []Grant{{"*", "*"}, {"tenants", "manage"}}) {
		t.Fatal("*/* covers everything")
	}
	if !CoversAll(rbacOnly, nil) {
		t.Fatal("an empty role is always grantable")
	}
}

func TestIsEffectiveAdministrator(t *testing.T) {
	if !IsEffectiveAdministrator([]Grant{{"*", "*"}}) {
		t.Fatal("*/* is administrator")
	}
	if !IsEffectiveAdministrator([]Grant{{"roles", "manage"}, {"permissions", "*"}, {"rbac", "manage"}}) {
		t.Fatal("the three management grants make an administrator")
	}
	if IsEffectiveAdministrator([]Grant{{"roles", "manage"}, {"permissions", "manage"}}) {
		t.Fatal("without rbac/manage it is not an administrator")
	}
}

func TestIsSensitiveGrant(t *testing.T) {
	for _, g := range []Grant{{"*", "read"}, {"invoices", "*"}, {"tenants", "manage"}} {
		if !IsSensitiveGrant(g) {
			t.Errorf("%v should be sensitive", g)
		}
	}
	if IsSensitiveGrant(Grant{"invoices", "read"}) {
		t.Error("invoices/read is not sensitive")
	}
}

func TestHasPermission_ManageImpliesRead(t *testing.T) {
	perms := []*identity.Permission{{TenantId: "t", Resource: "roles", Action: "manage"}}
	if !HasPermission("t", perms, "roles", "read") {
		t.Fatal("manage must imply read")
	}
	if HasPermission("t", []*identity.Permission{{TenantId: "t", Resource: "roles", Action: "read"}}, "roles", "manage") {
		t.Fatal("read must not imply manage")
	}
}
