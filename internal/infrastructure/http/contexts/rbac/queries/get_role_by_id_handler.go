package queries

import (
	rbacqueries "kali-auth-context/internal/application/rbac/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetRoleByIdHandler struct {
	query *rbacqueries.GetRoleByIdQuery
}

func NewGetRoleByIdHandler(query *rbacqueries.GetRoleByIdQuery) *GetRoleByIdHandler {
	return &GetRoleByIdHandler{query: query}
}

func (h *GetRoleByIdHandler) Handle(c *fiber.Ctx) error {
	role, err := h.query.Handle(c.UserContext(), identity.RoleId(c.Params("roleId")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(role)
}
