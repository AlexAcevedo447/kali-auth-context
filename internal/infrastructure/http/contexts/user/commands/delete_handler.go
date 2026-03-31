package commands

import (
	usercommands "kali-auth-context/internal/application/user/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type DeleteHandler struct {
	command *usercommands.DeleteUserCommand
}

func NewDeleteHandler(command *usercommands.DeleteUserCommand) *DeleteHandler {
	return &DeleteHandler{command: command}
}

func (h *DeleteHandler) Handle(c *fiber.Ctx) error {
	err := h.command.Execute(c.UserContext(), &usercommands.DeleteUserDto{
		TenantId: identity.TenantId(c.Query("tenant_id")),
		UserId:   identity.UserId(c.Params("userId")),
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
