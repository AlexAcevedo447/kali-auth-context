package commands

import (
	tenantcommands "kali-auth-context/internal/application/tenant/commands"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type CreateHandler struct {
	command *tenantcommands.CreateTenantCommand
}

func NewCreateHandler(command *tenantcommands.CreateTenantCommand) *CreateHandler {
	return &CreateHandler{command: command}
}

func (h *CreateHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &tenantcommands.CreateTenantDto{Name: req.Name})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}
