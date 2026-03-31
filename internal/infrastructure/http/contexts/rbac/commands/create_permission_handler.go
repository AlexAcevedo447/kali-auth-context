package commands

import (
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type CreatePermissionHandler struct {
	command *rbaccommands.CreatePermissionCommand
}

func NewCreatePermissionHandler(command *rbaccommands.CreatePermissionCommand) *CreatePermissionHandler {
	return &CreatePermissionHandler{command: command}
}

func (h *CreatePermissionHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		TenantId string `json:"tenant_id"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &rbaccommands.CreatePermissionDto{
		TenantId: identity.TenantId(req.TenantId),
		Resource: req.Resource,
		Action:   req.Action,
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}
