package commands

import (
	"context"
	"kali-auth-context/internal/domain/policies"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type UpdatePermissionDto struct {
	Id       identity.PermissionId
	TenantId identity.TenantId
	Resource string
	Action   string
}

type UpdatePermissionCommand struct {
	repo ports.IUpdatePermissionCommandRepository
}

func NewUpdatePermissionCommand(repo ports.IUpdatePermissionCommandRepository) *UpdatePermissionCommand {
	return &UpdatePermissionCommand{repo: repo}
}

func (c *UpdatePermissionCommand) Execute(ctx context.Context, dto *UpdatePermissionDto) error {
	resource, action, err := policies.ValidatePermission(dto.Resource, dto.Action)
	if err != nil {
		return err
	}
	permission, err := identity.NewPermission(dto.Id, dto.TenantId, resource, action)
	if err != nil {
		return err
	}

	return c.repo.Update(ctx, permission)
}
