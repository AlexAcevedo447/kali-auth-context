package commands

import (
	usercommands "kali-auth-context/internal/application/user/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type CreateHandler struct {
	command *usercommands.CreateUserCommand
}

func NewCreateHandler(command *usercommands.CreateUserCommand) *CreateHandler {
	return &CreateHandler{command: command}
}

func (h *CreateHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		TenantId             string `json:"tenant_id"`
		IdentificationNumber string `json:"identification_number"`
		Username             string `json:"username"`
		Email                string `json:"email"`
		Password             string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &usercommands.CreateUserDto{
		TenantId:             identity.TenantId(req.TenantId),
		IdentificationNumber: req.IdentificationNumber,
		Username:             req.Username,
		Email:                req.Email,
		Password:             req.Password,
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}
