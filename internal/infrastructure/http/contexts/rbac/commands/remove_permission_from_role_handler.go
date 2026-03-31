package commands

import (
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type RemovePermissionFromRoleHandler struct {
	command *rbaccommands.RemovePermissionFromRoleCommand
}

func NewRemovePermissionFromRoleHandler(command *rbaccommands.RemovePermissionFromRoleCommand) *RemovePermissionFromRoleHandler {
	return &RemovePermissionFromRoleHandler{command: command}
}

func (h *RemovePermissionFromRoleHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		TenantId     string `json:"tenant_id"`
		RoleId       string `json:"role_id"`
		PermissionId string `json:"permission_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &rbaccommands.RemovePermissionFromRoleDto{
		TenantId:     identity.TenantId(req.TenantId),
		RoleId:       identity.RoleId(req.RoleId),
		PermissionId: identity.PermissionId(req.PermissionId),
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)

}
