package queries

import (
	tenantqueries "kali-auth-context/internal/application/tenant/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type GetByIdHandler struct {
	query *tenantqueries.GetTenantByIdQuery
}

func NewGetByIdHandler(query *tenantqueries.GetTenantByIdQuery) *GetByIdHandler {
	return &GetByIdHandler{query: query}
}

func (h *GetByIdHandler) Handle(c *fiber.Ctx) error {
	tenant, err := h.query.Handle(c.UserContext(), identity.TenantId(c.Params("tenantId")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(tenant)
}
