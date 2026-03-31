package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type RemovePermissionFromRoleDto struct {
	TenantId     identity.TenantId
	RoleId       identity.RoleId
	PermissionId identity.PermissionId
}

type RemovePermissionFromRoleCommand struct {
	repo ports.IRemovePermissionFromRoleCommandRepository
}

func NewRemovePermissionFromRoleCommand(repo ports.IRemovePermissionFromRoleCommandRepository) *RemovePermissionFromRoleCommand {
	return &RemovePermissionFromRoleCommand{repo: repo}
}

func (c *RemovePermissionFromRoleCommand) Execute(ctx context.Context, dto *RemovePermissionFromRoleDto) error {
	if dto.TenantId == "" {
		return identity.ErrTenantRequired
	}
	if dto.RoleId == "" {
		return identity.ErrRoleIdRequired
	}
	if dto.PermissionId == "" {
		return identity.ErrPermissionIdRequired
	}

	return c.repo.Remove(ctx, dto.TenantId, dto.RoleId, dto.PermissionId)
}
