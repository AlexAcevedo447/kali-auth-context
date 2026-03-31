package commands

import (
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type RemoveRoleFromUserHandler struct {
	command *rbaccommands.RemoveRoleFromUserCommand
}

func NewRemoveRoleFromUserHandler(command *rbaccommands.RemoveRoleFromUserCommand) *RemoveRoleFromUserHandler {
	return &RemoveRoleFromUserHandler{command: command}
}

func (h *RemoveRoleFromUserHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		TenantId string `json:"tenant_id"`
		UserId   string `json:"user_id"`
		RoleId   string `json:"role_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &rbaccommands.RemoveRoleFromUserDto{
		TenantId: identity.TenantId(req.TenantId),
		UserId:   identity.UserId(req.UserId),
		RoleId:   identity.RoleId(req.RoleId),
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
