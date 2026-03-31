package commands

import (
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type UpdatePermissionHandler struct {
	command *rbaccommands.UpdatePermissionCommand
}

func NewUpdatePermissionHandler(command *rbaccommands.UpdatePermissionCommand) *UpdatePermissionHandler {
	return &UpdatePermissionHandler{command: command}
}

func (h *UpdatePermissionHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		TenantId string `json:"tenant_id"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &rbaccommands.UpdatePermissionDto{
		Id:       identity.PermissionId(c.Params("permissionId")),
		TenantId: identity.TenantId(req.TenantId),
		Resource: req.Resource,
		Action:   req.Action,
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
