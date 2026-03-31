package queries

import (
	userqueries "kali-auth-context/internal/application/user/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetByIdHandler struct {
	query *userqueries.GetUserByIdQuery
}

func NewGetByIdHandler(query *userqueries.GetUserByIdQuery) *GetByIdHandler {
	return &GetByIdHandler{query: query}
}

func (h *GetByIdHandler) Handle(c *fiber.Ctx) error {
	user, err := h.query.Handle(identity.TenantId(c.Query("tenant_id")), identity.UserId(c.Params("userId")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
