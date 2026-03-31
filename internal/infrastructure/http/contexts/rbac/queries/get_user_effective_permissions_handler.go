package queries

import (
	rbacqueries "kali-auth-context/internal/application/rbac/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetUserEffectivePermissionsHandler struct {
	query *rbacqueries.GetUserEffectivePermissionsQuery
}

func NewGetUserEffectivePermissionsHandler(query *rbacqueries.GetUserEffectivePermissionsQuery) *GetUserEffectivePermissionsHandler {
	return &GetUserEffectivePermissionsHandler{query: query}
}

func (h *GetUserEffectivePermissionsHandler) Handle(c *fiber.Ctx) error {
	permissions, err := h.query.Handle(
		c.UserContext(),
		identity.TenantId(c.Query("tenant_id")),
		identity.UserId(c.Query("user_id")),
	)
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(permissions)
}
