package middleware

import "testing"

// Simulación de "no dejar a la empresa sin administrador efectivo".
func TestCountAdministrators_Simulation(t *testing.T) {
	row := func(user, role, perm, res, act string) assignmentRow {
		return assignmentRow{UserID: user, RoleID: role, PermissionID: perm, Resource: res, Action: act}
	}
	rows := []assignmentRow{
		row("ana", "admin", "p-all", "*", "*"),
		row("beto", "soporte", "p-roles", "roles", "manage"),
		row("beto", "soporte", "p-perms", "permissions", "manage"),
		row("beto", "soporte", "p-rbac", "rbac", "manage"),
		row("caro", "factura", "p-inv", "invoices", "manage"),
	}
	if got := countAdministrators(rows, nil); got != 2 {
		t.Fatalf("expected 2 administrators, got %d", got)
	}
	// Quitar el rol admin a ana: queda beto.
	if got := countAdministrators(rows, func(r assignmentRow) bool { return r.UserID == "ana" && r.RoleID == "admin" }); got != 1 {
		t.Fatalf("expected 1 administrator, got %d", got)
	}
	// Además quitar rbac/manage al rol soporte: no queda ninguno.
	skip := func(r assignmentRow) bool {
		return (r.UserID == "ana" && r.RoleID == "admin") || (r.RoleID == "soporte" && r.PermissionID == "p-rbac")
	}
	if got := countAdministrators(rows, skip); got != 0 {
		t.Fatalf("expected 0 administrators, got %d", got)
	}
}
