package queries

import (
	tenantqueries "kali-auth-context/internal/application/tenant/queries"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type ListHandler struct {
	query *tenantqueries.ListTenantsQuery
}

func NewListHandler(query *tenantqueries.ListTenantsQuery) *ListHandler {
	return &ListHandler{query: query}
}

func (h *ListHandler) Handle(c *fiber.Ctx) error {
	tenants, err := h.query.Handle(c.UserContext())
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(tenants)
}
