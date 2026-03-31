package commands

import (
	tenantcommands "kali-auth-context/internal/application/tenant/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type SuspendHandler struct {
	command *tenantcommands.SuspendTenantCommand
}

func NewSuspendHandler(command *tenantcommands.SuspendTenantCommand) *SuspendHandler {
	return &SuspendHandler{command: command}
}

func (h *SuspendHandler) Handle(c *fiber.Ctx) error {
	err := h.command.Execute(c.UserContext(), identity.TenantId(c.Params("tenantId")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
