package commands

import (
	tenantcommands "kali-auth-context/internal/application/tenant/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type ActivateHandler struct {
	command *tenantcommands.ActivateTenantCommand
}

func NewActivateHandler(command *tenantcommands.ActivateTenantCommand) *ActivateHandler {
	return &ActivateHandler{command: command}
}

func (h *ActivateHandler) Handle(c *fiber.Ctx) error {
	err := h.command.Execute(c.UserContext(), identity.TenantId(c.Params("tenantId")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
