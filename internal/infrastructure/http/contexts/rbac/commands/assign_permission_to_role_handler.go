package commands

import (
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type AssignPermissionToRoleHandler struct {
	command *rbaccommands.AssignPermissionToRoleCommand
}

func NewAssignPermissionToRoleHandler(command *rbaccommands.AssignPermissionToRoleCommand) *AssignPermissionToRoleHandler {
	return &AssignPermissionToRoleHandler{command: command}
}

func (h *AssignPermissionToRoleHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		TenantId     string `json:"tenant_id"`
		RoleId       string `json:"role_id"`
		PermissionId string `json:"permission_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &rbaccommands.AssignPermissionToRoleDto{
		TenantId:     identity.TenantId(req.TenantId),
		RoleId:       identity.RoleId(req.RoleId),
		PermissionId: identity.PermissionId(req.PermissionId),
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
