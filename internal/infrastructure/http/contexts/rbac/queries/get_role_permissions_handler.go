package queries

import (
	rbacqueries "kali-auth-context/internal/application/rbac/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetRolePermissionsHandler struct {
	query *rbacqueries.GetRolePermissionsQuery
}

func NewGetRolePermissionsHandler(query *rbacqueries.GetRolePermissionsQuery) *GetRolePermissionsHandler {
	return &GetRolePermissionsHandler{query: query}
}

func (h *GetRolePermissionsHandler) Handle(c *fiber.Ctx) error {
	permissions, err := h.query.Handle(
		c.UserContext(),
		identity.TenantId(c.Query("tenant_id")),
		identity.RoleId(c.Query("role_id")),
	)
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(permissions)
}
