package queries

import (
	tenantqueries "kali-auth-context/internal/application/tenant/queries"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetByNameHandler struct {
	query *tenantqueries.GetTenantByNameQuery
}

func NewGetByNameHandler(query *tenantqueries.GetTenantByNameQuery) *GetByNameHandler {
	return &GetByNameHandler{query: query}
}

func (h *GetByNameHandler) Handle(c *fiber.Ctx) error {
	tenant, err := h.query.Handle(c.UserContext(), c.Query("name"))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(tenant)
}
