package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

const localsEffectiveTenant = "tenant.effective"

// TenantGuards aplica el aislamiento multi-tenant sobre las rutas protegidas.
//
// Regla: el tenant efectivo de cada petición es el del token (claim "tid"). El
// tenant_id que envíe el cliente nunca es autoridad. La única excepción es el
// tenant de plataforma (el tenant maestro sembrado), que administra empresas y
// puede indicar un tenant destino existente.
type TenantGuards struct {
	db               *bun.DB
	platformTenantID string
	authorize        Authorizer
}

// Authorizer evalúa los permisos efectivos (RBAC) de un usuario. En producción
// es AuthorizeQuery (roles → permisos con comodines).
type Authorizer func(ctx context.Context, tenantID, userID, resource, action string) (bool, error)

// Permiso RBAC que habilita la administración de empresas desde la empresa de
// plataforma (el seed maestro lo cubre con "*"/"*").
const (
	PlatformAdminResource = "tenants"
	PlatformAdminAction   = "manage"
)

func NewTenantGuards(db *bun.DB, platformTenantID string) *TenantGuards {
	return &TenantGuards{db: db, platformTenantID: platformTenantID}
}

// WithAuthorizer conecta la evaluación RBAC. Sin ella, ninguna operación de
// administración de plataforma se autoriza (falla cerrado).
func (g *TenantGuards) WithAuthorizer(a Authorizer) *TenantGuards {
	g.authorize = a
	return g
}

func (g *TenantGuards) isPlatform(c *fiber.Ctx) bool {
	tid, _ := c.Locals("auth.tid").(string)
	return g.platformTenantID != "" && tid == g.platformTenantID
}

// isPlatformAdmin: usuario de la empresa de plataforma CON el permiso
// tenants/manage. Ser de la empresa de plataforma no basta.
func (g *TenantGuards) isPlatformAdmin(c *fiber.Ctx) (bool, error) {
	if !g.isPlatform(c) || g.authorize == nil {
		return false, nil
	}
	tid, _ := c.Locals("auth.tid").(string)
	sub, _ := c.Locals("auth.sub").(string)
	if sub == "" {
		return false, nil
	}
	return g.authorize(c.UserContext(), tid, sub, PlatformAdminResource, PlatformAdminAction)
}

// ActiveTenant rechaza cualquier operación protegida si la empresa del token no
// está activa. Se consulta en cada petición, así que suspender una empresa
// bloquea también los tokens ya emitidos; al reactivarla, los tokens aún
// vigentes vuelven a funcionar (la vigencia la fija JWT_ACCESS_TOKEN_TTL_MINUTES).
func (g *TenantGuards) ActiveTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tid, _ := c.Locals("auth.tid").(string)
		if tid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token is missing required claims"})
		}
		var status string
		err := g.db.NewSelect().Table("tenants").Column("status").Where("id = ?", tid).Scan(c.UserContext(), &status)
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid access token"})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		if status != "ACTIVE" {
			return forbidden(c, "tenant is suspended")
		}
		// El usuario del token debe seguir existiendo en su empresa: un usuario
		// eliminado pierde el acceso también con los tokens ya emitidos.
		if sub, _ := c.Locals("auth.sub").(string); sub != "" {
			exists, err := g.db.NewSelect().TableExpr("user_models").
				Where("id = ? AND tenant_id = ?", sub, tid).Exists(c.UserContext())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
			}
			if !exists {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid access token"})
			}
		}
		return c.Next()
	}
}

func forbidden(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": msg})
}

func notFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "resource not found"})
}

// Scope fija el tenant de la petición al del token, tanto en el query como en
// el cuerpo JSON. Si el cliente envía un tenant_id distinto se rechaza con 403
// (salvo el tenant de plataforma, que puede apuntar a un tenant existente).
func (g *TenantGuards) Scope() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tid, _ := c.Locals("auth.tid").(string)
		if tid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token is missing required claims"})
		}
		platform := g.isPlatform(c)
		effective := tid

		resolve := func(requested string) (string, error) {
			if requested == "" || requested == tid {
				return tid, nil
			}
			// Apuntar a otra empresa es administración de plataforma: exige el
			// permiso tenants/manage, no solo pertenecer a la empresa de plataforma.
			if !platform {
				return "", errTenantMismatch
			}
			admin, err := g.isPlatformAdmin(c)
			if err != nil {
				return "", err
			}
			if !admin {
				return "", errTenantMismatch
			}
			exists, err := g.exists(c.UserContext(), "tenants", requested)
			if err != nil {
				return "", err
			}
			if !exists {
				return "", errTenantNotFound
			}
			return requested, nil
		}

		// Query string (se conserva el valor original del cliente para comparar con el body)
		queryTenant := c.Query("tenant_id")
		requested, err := resolve(queryTenant)
		if err != nil {
			return g.writeScopeError(c, err)
		}
		effective = requested
		c.Request().URI().QueryArgs().Set("tenant_id", effective)

		// Cuerpo JSON
		body := c.Body()
		if len(body) > 0 && strings.Contains(strings.ToLower(c.Get(fiber.HeaderContentType)), "json") {
			var payload map[string]any
			if err := json.Unmarshal(body, &payload); err == nil {
				bodyTenant, _ := payload["tenant_id"].(string)
				resolved, err := resolve(bodyTenant)
				if err != nil {
					return g.writeScopeError(c, err)
				}
				if bodyTenant == "" {
					resolved = effective
				} else if queryTenant != "" && resolved != effective {
					return forbidden(c, "tenant mismatch")
				}
				effective = resolved
				payload["tenant_id"] = effective
				rewritten, err := json.Marshal(payload)
				if err != nil {
					return err
				}
				c.Request().SetBody(rewritten)
				c.Request().URI().QueryArgs().Set("tenant_id", effective)
			}
		}

		c.Locals(localsEffectiveTenant, effective)
		return c.Next()
	}
}

var (
	errTenantMismatch = errors.New("tenant mismatch")
	errTenantNotFound = errors.New("tenant not found")
)

func (g *TenantGuards) writeScopeError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, errTenantMismatch):
		return forbidden(c, "tenant mismatch")
	case errors.Is(err, errTenantNotFound):
		return notFound(c)
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}

// PlatformAdmin restringe una ruta al administrador de plataforma: usuario de
// la empresa de plataforma con el permiso tenants/manage (RBAC). Cualquier otro
// usuario, incluido uno de la plataforma sin ese permiso, recibe 403.
func (g *TenantGuards) PlatformAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		admin, err := g.isPlatformAdmin(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		if !admin {
			return forbidden(c, "platform administrator required")
		}
		return c.Next()
	}
}

// OwnedRef describe un identificador de la petición que debe pertenecer al
// tenant efectivo: Source es "param", "query" o "body".
type OwnedRef struct {
	Source string
	Field  string
	Table  string
}

// Owned verifica que cada recurso referenciado pertenezca al tenant efectivo.
// Un recurso de otra empresa se responde como inexistente (404), sin revelar
// que existe.
func (g *TenantGuards) Owned(refs ...OwnedRef) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenant, _ := c.Locals(localsEffectiveTenant).(string)
		var payload map[string]any
		for _, ref := range refs {
			var id string
			switch ref.Source {
			case "param":
				id = c.Params(ref.Field)
			case "query":
				id = c.Query(ref.Field)
			case "body":
				if payload == nil {
					_ = json.Unmarshal(c.Body(), &payload)
				}
				id, _ = payload[ref.Field].(string)
			}
			if id == "" {
				continue // los handlers validan los campos obligatorios
			}
			owner, err := g.ownerOf(c.UserContext(), ref.Table, id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return notFound(c)
				}
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
			}
			if owner != tenant {
				return notFound(c)
			}
		}
		return c.Next()
	}
}

// allowedTables evita construir SQL con nombres de tabla arbitrarios.
var allowedTables = map[string]bool{"roles": true, "permissions": true, "user_models": true, "tenants": true}

func (g *TenantGuards) ownerOf(ctx context.Context, table, id string) (string, error) {
	if !allowedTables[table] || table == "tenants" {
		return "", errors.New("unsupported table")
	}
	var owner sql.NullString
	err := g.db.NewSelect().TableExpr(table).Column("tenant_id").Where("id = ?", id).Limit(1).Scan(ctx, &owner)
	if err != nil {
		return "", err
	}
	return owner.String, nil
}

func (g *TenantGuards) exists(ctx context.Context, table, id string) (bool, error) {
	if !allowedTables[table] {
		return false, errors.New("unsupported table")
	}
	return g.db.NewSelect().TableExpr(table).Where("id = ?", id).Exists(ctx)
}
