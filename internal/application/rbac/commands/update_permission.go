package commands

import (
	"context"

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
	permission, err := identity.NewPermission(dto.Id, dto.TenantId, dto.Resource, dto.Action)
	if err != nil {
		return err
	}

	return c.repo.Update(ctx, permission)
}
