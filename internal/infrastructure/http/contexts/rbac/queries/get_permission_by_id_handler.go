package queries

import (
	rbacqueries "kali-auth-context/internal/application/rbac/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetPermissionByIdHandler struct {
	query *rbacqueries.GetPermissionByIdQuery
}

func NewGetPermissionByIdHandler(query *rbacqueries.GetPermissionByIdQuery) *GetPermissionByIdHandler {
	return &GetPermissionByIdHandler{query: query}
}

func (h *GetPermissionByIdHandler) Handle(c *fiber.Ctx) error {
	permission, err := h.query.Handle(c.UserContext(), identity.PermissionId(c.Params("permissionId")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(permission)
}
