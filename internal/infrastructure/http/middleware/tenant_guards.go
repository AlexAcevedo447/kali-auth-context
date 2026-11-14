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
}

func NewTenantGuards(db *bun.DB, platformTenantID string) *TenantGuards {
	return &TenantGuards{db: db, platformTenantID: platformTenantID}
}

func (g *TenantGuards) isPlatform(c *fiber.Ctx) bool {
	tid, _ := c.Locals("auth.tid").(string)
	return g.platformTenantID != "" && tid == g.platformTenantID
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
			if !platform {
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

// PlatformOnly restringe la administración de empresas (tenants) al tenant de plataforma.
func (g *TenantGuards) PlatformOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !g.isPlatform(c) {
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
