package commands

import (
	authcommands "kali-auth-context/internal/application/auth/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type LoginHandler struct {
	command *authcommands.LoginCommand
}

func NewLoginHandler(command *authcommands.LoginCommand) *LoginHandler {
	return &LoginHandler{command: command}
}

type loginRequest struct {
	TenantId string `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *LoginHandler) Handle(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := h.command.Execute(&authcommands.LoginDto{
		TenantId: identity.TenantId(req.TenantId),
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
