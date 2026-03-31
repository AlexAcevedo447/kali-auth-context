package commands

import (
	tenantcommands "kali-auth-context/internal/application/tenant/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type UpdateHandler struct {
	command *tenantcommands.UpdateTenantCommand
}

func NewUpdateHandler(command *tenantcommands.UpdateTenantCommand) *UpdateHandler {
	return &UpdateHandler{command: command}
}

func (h *UpdateHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &tenantcommands.UpdateTenantDto{
		Id:     identity.TenantId(c.Params("tenantId")),
		Name:   req.Name,
		Status: identity.TenantStatus(req.Status),
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
