package commands

import (
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type DeletePermissionHandler struct {
	command *rbaccommands.DeletePermissionCommand
}

func NewDeletePermissionHandler(command *rbaccommands.DeletePermissionCommand) *DeletePermissionHandler {
	return &DeletePermissionHandler{command: command}
}

func (h *DeletePermissionHandler) Handle(c *fiber.Ctx) error {
	err := h.command.Execute(c.UserContext(), identity.PermissionId(c.Params("permissionId")))
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
