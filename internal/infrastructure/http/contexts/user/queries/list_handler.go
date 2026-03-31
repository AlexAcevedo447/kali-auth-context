package queries

import (
	userqueries "kali-auth-context/internal/application/user/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type ListHandler struct {
	query *userqueries.ListUsersQuery
}

func NewListHandler(query *userqueries.ListUsersQuery) *ListHandler {
	return &ListHandler{query: query}
}

func (h *ListHandler) Handle(c *fiber.Ctx) error {
	users, err := h.query.Handle(identity.TenantId(c.Query("tenant_id")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
