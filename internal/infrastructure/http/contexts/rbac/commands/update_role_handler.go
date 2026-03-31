package commands

import (
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type UpdateRoleHandler struct {
	command *rbaccommands.UpdateRoleCommand
}

func NewUpdateRoleHandler(command *rbaccommands.UpdateRoleCommand) *UpdateRoleHandler {
	return &UpdateRoleHandler{command: command}
}

func (h *UpdateRoleHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		TenantId    string `json:"tenant_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := h.command.Execute(c.UserContext(), &rbaccommands.UpdateRoleDto{
		Id:          identity.RoleId(c.Params("roleId")),
		TenantId:    identity.TenantId(req.TenantId),
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return shared.WriteError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
