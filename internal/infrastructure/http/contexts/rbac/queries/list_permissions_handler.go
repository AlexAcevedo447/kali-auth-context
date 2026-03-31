package queries

import (
	rbacqueries "kali-auth-context/internal/application/rbac/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type ListPermissionsHandler struct {
	query *rbacqueries.ListPermissionsQuery
}

func NewListPermissionsHandler(query *rbacqueries.ListPermissionsQuery) *ListPermissionsHandler {
	return &ListPermissionsHandler{query: query}
}

func (h *ListPermissionsHandler) Handle(c *fiber.Ctx) error {
	permissions, err := h.query.Handle(c.UserContext(), identity.TenantId(c.Query("tenant_id")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(permissions)
}
