package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"
)

// RBACGuards aplica la autorización por operación del incremento 2 (ver
// MATRIZ_PERMISOS.md en invoice-kali-front): permiso exigido por operación,
// no escalamiento de privilegios, protección del administrador sembrado y de
// que la empresa conserve al menos un administrador efectivo.
//
// Se monta después de jwtAuth, ActiveTenant, Scope y Owned: el tenant efectivo
// y la pertenencia de los ids referenciados ya están verificados.
type RBACGuards struct {
	db               *bun.DB
	platformTenantID string
}

func NewRBACGuards(db *bun.DB, platformTenantID string) *RBACGuards {
	return &RBACGuards{db: db, platformTenantID: platformTenantID}
}

// --- Lectura de datos RBAC -------------------------------------------------------

type assignmentRow struct {
	UserID       string `bun:"user_id"`
	RoleID       string `bun:"role_id"`
	PermissionID string `bun:"permission_id"`
	Resource     string `bun:"resource"`
	Action       string `bun:"action"`
}

// EffectiveGrants: permisos efectivos (roles → permisos) de un usuario en su
// empresa. Las asignaciones huérfanas no otorgan nada.
func (g *RBACGuards) EffectiveGrants(ctx context.Context, tenantID, userID string) ([]policies.Grant, error) {
	var rows []assignmentRow
	err := g.db.NewRaw(`SELECT ur.user_id, ur.role_id, rp.permission_id, p.resource, p.action
		FROM user_roles ur
		JOIN role_permissions rp ON rp.role_id = ur.role_id
		JOIN permissions p ON p.id = rp.permission_id AND p.tenant_id = ur.tenant_id
		WHERE ur.tenant_id = ? AND ur.user_id = ?`, tenantID, userID).Scan(ctx, &rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return grantsOf(rows, nil), nil
}

func (g *RBACGuards) rolePermissionGrants(ctx context.Context, roleID string) ([]policies.Grant, error) {
	var rows []assignmentRow
	err := g.db.NewRaw(`SELECT '' AS user_id, rp.role_id, rp.permission_id, p.resource, p.action
		FROM role_permissions rp JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = ?`, roleID).Scan(ctx, &rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return grantsOf(rows, nil), nil
}

func (g *RBACGuards) permissionGrant(ctx context.Context, permissionID string) (tenantID string, grant policies.Grant, err error) {
	var row struct {
		TenantID string `bun:"tenant_id"`
		Resource string `bun:"resource"`
		Action   string `bun:"action"`
	}
	err = g.db.NewRaw(`SELECT tenant_id, resource, action FROM permissions WHERE id = ?`, permissionID).Scan(ctx, &row)
	return row.TenantID, policies.Grant{Resource: row.Resource, Action: row.Action}, err
}

func (g *RBACGuards) role(ctx context.Context, roleID string) (tenantID, name string, err error) {
	var row struct {
		TenantID string `bun:"tenant_id"`
		Name     string `bun:"name"`
	}
	err = g.db.NewRaw(`SELECT tenant_id, name FROM roles WHERE id = ?`, roleID).Scan(ctx, &row)
	return row.TenantID, row.Name, err
}

func (g *RBACGuards) tenantAssignments(ctx context.Context, tenantID string) ([]assignmentRow, error) {
	var rows []assignmentRow
	err := g.db.NewRaw(`SELECT ur.user_id, ur.role_id, rp.permission_id, p.resource, p.action
		FROM user_roles ur
		JOIN role_permissions rp ON rp.role_id = ur.role_id
		JOIN permissions p ON p.id = rp.permission_id AND p.tenant_id = ur.tenant_id
		WHERE ur.tenant_id = ?`, tenantID).Scan(ctx, &rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return rows, nil
}

func grantsOf(rows []assignmentRow, skip func(assignmentRow) bool) []policies.Grant {
	out := make([]policies.Grant, 0, len(rows))
	for _, r := range rows {
		if skip != nil && skip(r) {
			continue
		}
		out = append(out, policies.Grant{Resource: r.Resource, Action: r.Action})
	}
	return out
}

// countAdministrators cuenta usuarios de la empresa que son administradores
// efectivos, ignorando las filas que skip descarte (simulación de un cambio).
func countAdministrators(rows []assignmentRow, skip func(assignmentRow) bool) int {
	byUser := map[string][]policies.Grant{}
	for _, r := range rows {
		if skip != nil && skip(r) {
			continue
		}
		byUser[r.UserID] = append(byUser[r.UserID], policies.Grant{Resource: r.Resource, Action: r.Action})
	}
	n := 0
	for _, grants := range byUser {
		if policies.IsEffectiveAdministrator(grants) {
			n++
		}
	}
	return n
}

// --- Utilidades de petición --------------------------------------------------------

func actor(c *fiber.Ctx) (tenantID, userID string) {
	tenantID, _ = c.Locals("auth.tid").(string)
	userID, _ = c.Locals("auth.sub").(string)
	return
}

func effectiveTenant(c *fiber.Ctx) string {
	if t, ok := c.Locals(localsEffectiveTenant).(string); ok && t != "" {
		return t
	}
	t, _ := c.Locals("auth.tid").(string)
	return t
}

func requestValue(c *fiber.Ctx, source, field string) string {
	switch source {
	case "param":
		return c.Params(field)
	case "query":
		return c.Query(field)
	case "body":
		var payload map[string]any
		if err := json.Unmarshal(c.Body(), &payload); err != nil {
			return ""
		}
		v, _ := payload[field].(string)
		return v
	}
	return ""
}

func writeErr(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}

func internalError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}

func (g *RBACGuards) actorGrants(c *fiber.Ctx) ([]policies.Grant, error) {
	tid, sub := actor(c)
	if tid == "" || sub == "" {
		return nil, nil
	}
	return g.EffectiveGrants(c.UserContext(), tid, sub)
}

// --- Guardas -------------------------------------------------------------------------

// Require exige el permiso (recurso, acción) con los comodines y manage⇒read.
func (g *RBACGuards) Require(resource, action string) fiber.Handler {
	want := policies.Grant{Resource: resource, Action: action}
	return func(c *fiber.Ctx) error {
		grants, err := g.actorGrants(c)
		if err != nil {
			return internalError(c)
		}
		if !policies.Covers(grants, want) {
			return writeErr(c, fiber.StatusForbidden, identity.ErrInsufficientPermissions)
		}
		return c.Next()
	}
}

// SelfOr permite consultar los datos propios (el id indicado es el del token y
// la empresa efectiva es la suya); para otro usuario exige (recurso, acción).
func (g *RBACGuards) SelfOr(source, field, resource, action string) fiber.Handler {
	require := g.Require(resource, action)
	return func(c *fiber.Ctx) error {
		tid, sub := actor(c)
		target := requestValue(c, source, field)
		if target != "" && target == sub && effectiveTenant(c) == tid {
			return c.Next()
		}
		return require(c)
	}
}

// NoEscalationUserRole: para asignar o quitar un rol, quien opera debe cubrir
// todos los permisos de ese rol.
func (g *RBACGuards) NoEscalationUserRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleGrants, err := g.rolePermissionGrants(c.UserContext(), requestValue(c, "body", "role_id"))
		if err != nil {
			return internalError(c)
		}
		mine, err := g.actorGrants(c)
		if err != nil {
			return internalError(c)
		}
		if !policies.CoversAll(mine, roleGrants) {
			return writeErr(c, fiber.StatusForbidden, identity.ErrPrivilegeEscalation)
		}
		return c.Next()
	}
}

// NoEscalationRolePermission: para asignar o quitar un permiso a un rol, quien
// opera debe cubrir ese permiso.
func (g *RBACGuards) NoEscalationRolePermission() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, grant, err := g.permissionGrant(c.UserContext(), requestValue(c, "body", "permission_id"))
		if errors.Is(err, sql.ErrNoRows) {
			return notFound(c)
		}
		if err != nil {
			return internalError(c)
		}
		mine, err := g.actorGrants(c)
		if err != nil {
			return internalError(c)
		}
		if !policies.Covers(mine, grant) {
			return writeErr(c, fiber.StatusForbidden, identity.ErrPrivilegeEscalation)
		}
		return c.Next()
	}
}

// NoEscalationPermissionWrite: crear o editar un permiso con comodín o de
// administración de plataforma exige cubrirlo.
func (g *RBACGuards) NoEscalationPermissionWrite() fiber.Handler {
	return func(c *fiber.Ctx) error {
		want := policies.Grant{Resource: requestValue(c, "body", "resource"), Action: requestValue(c, "body", "action")}
		if !policies.IsSensitiveGrant(want) {
			return c.Next()
		}
		mine, err := g.actorGrants(c)
		if err != nil {
			return internalError(c)
		}
		if !policies.Covers(mine, want) {
			return writeErr(c, fiber.StatusForbidden, identity.ErrPrivilegeEscalation)
		}
		return c.Next()
	}
}

func (g *RBACGuards) isSeedRole(tenantID, name string) bool {
	return g.platformTenantID != "" && tenantID == g.platformTenantID &&
		strings.ToLower(strings.TrimSpace(name)) == "admin"
}

// ProtectSeedRole impide modificar o eliminar el rol "admin" sembrado de la
// plataforma y quitarle permisos.
func (g *RBACGuards) ProtectSeedRole(source, field string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID, name, err := g.role(c.UserContext(), requestValue(c, source, field))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return internalError(c)
		}
		if err == nil && g.isSeedRole(tenantID, name) {
			return writeErr(c, fiber.StatusConflict, identity.ErrProtectedRole)
		}
		return c.Next()
	}
}

// ProtectSeedPermission impide modificar o eliminar el permiso */* de la plataforma.
func (g *RBACGuards) ProtectSeedPermission(source, field string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID, grant, err := g.permissionGrant(c.UserContext(), requestValue(c, source, field))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return internalError(c)
		}
		r := policies.NormalizePermissionValue(grant.Resource)
		a := policies.NormalizePermissionValue(grant.Action)
		if err == nil && g.platformTenantID != "" && tenantID == g.platformTenantID && r == policies.Wildcard && a == policies.Wildcard {
			return writeErr(c, fiber.StatusConflict, identity.ErrProtectedPermission)
		}
		return c.Next()
	}
}

// Tipos de cambio que pueden dejar a una empresa sin administrador efectivo.
const (
	ChangeRemoveUserRole       = "remove-user-role"
	ChangeRemoveRolePermission = "remove-role-permission"
)

// KeepAdministrator bloquea con 409 un cambio que deje a la empresa sin
// administrador efectivo (solo si antes había al menos uno).
func (g *RBACGuards) KeepAdministrator(change string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := g.tenantAssignments(c.UserContext(), effectiveTenant(c))
		if err != nil {
			return internalError(c)
		}
		var skip func(assignmentRow) bool
		switch change {
		case ChangeRemoveUserRole:
			user, role := requestValue(c, "body", "user_id"), requestValue(c, "body", "role_id")
			skip = func(r assignmentRow) bool { return r.UserID == user && r.RoleID == role }
		case ChangeRemoveRolePermission:
			role, perm := requestValue(c, "body", "role_id"), requestValue(c, "body", "permission_id")
			skip = func(r assignmentRow) bool { return r.RoleID == role && r.PermissionID == perm }
		default:
			return internalError(c)
		}
		if countAdministrators(rows, nil) > 0 && countAdministrators(rows, skip) == 0 {
			return writeErr(c, fiber.StatusConflict, identity.ErrLastAdministrator)
		}
		return c.Next()
	}
}

// --- Handlers ------------------------------------------------------------------------

// SessionHandler: GET /auth/session. Solo llega aquí si el token es válido y su
// empresa está activa (ActiveTenant). Devuelve los permisos efectivos propios.
func (g *RBACGuards) SessionHandler(c *fiber.Ctx) error {
	tid, sub := actor(c)
	grants, err := g.EffectiveGrants(c.UserContext(), tid, sub)
	if err != nil {
		return internalError(c)
	}
	normalized := make([]policies.Grant, 0, len(grants))
	seen := map[policies.Grant]bool{}
	for _, gr := range grants {
		n := policies.Grant{Resource: policies.NormalizePermissionValue(gr.Resource), Action: policies.NormalizePermissionValue(gr.Action)}
		if !seen[n] {
			seen[n] = true
			normalized = append(normalized, n)
		}
	}
	return c.JSON(fiber.Map{"tenant_id": tid, "user_id": sub, "status": "ACTIVE", "permissions": normalized})
}

// CatalogHandler: GET /permissions/catalog.
func CatalogHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"resources": policies.CatalogEntries()})
}
