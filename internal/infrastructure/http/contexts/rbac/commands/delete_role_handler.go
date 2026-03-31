package commands

import (
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type DeleteRoleHandler struct {
	command *rbaccommands.DeleteRoleCommand
}

func NewDeleteRoleHandler(command *rbaccommands.DeleteRoleCommand) *DeleteRoleHandler {
	return &DeleteRoleHandler{command: command}
}

func (h *DeleteRoleHandler) Handle(c *fiber.Ctx) error {
	err := h.command.Execute(c.UserContext(), identity.RoleId(c.Params("roleId")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
