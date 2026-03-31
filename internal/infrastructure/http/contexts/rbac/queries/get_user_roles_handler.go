package queries

import (
	rbacqueries "kali-auth-context/internal/application/rbac/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetUserRolesHandler struct {
	query *rbacqueries.GetUserRolesQuery
}

func NewGetUserRolesHandler(query *rbacqueries.GetUserRolesQuery) *GetUserRolesHandler {
	return &GetUserRolesHandler{query: query}
}

func (h *GetUserRolesHandler) Handle(c *fiber.Ctx) error {
	roles, err := h.query.Handle(
		c.UserContext(),
		identity.TenantId(c.Query("tenant_id")),
		identity.UserId(c.Query("user_id")),
	)
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(roles)
}
