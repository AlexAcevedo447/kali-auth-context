package queries

import (
	userqueries "kali-auth-context/internal/application/user/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetByEmailHandler struct {
	query *userqueries.GetUserByEmailQuery
}

func NewGetByEmailHandler(query *userqueries.GetUserByEmailQuery) *GetByEmailHandler {
	return &GetByEmailHandler{query: query}
}

func (h *GetByEmailHandler) Handle(c *fiber.Ctx) error {
	user, err := h.query.Handle(identity.TenantId(c.Query("tenant_id")), c.Query("email"))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
