package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type AssignPermissionToRoleDto struct {
	TenantId     identity.TenantId
	RoleId       identity.RoleId
	PermissionId identity.PermissionId
}

type AssignPermissionToRoleCommand struct {
	repo ports.IAssignPermissionToRoleCommandRepository
}

func NewAssignPermissionToRoleCommand(repo ports.IAssignPermissionToRoleCommandRepository) *AssignPermissionToRoleCommand {
	return &AssignPermissionToRoleCommand{repo: repo}
}

func (c *AssignPermissionToRoleCommand) Execute(ctx context.Context, dto *AssignPermissionToRoleDto) error {
	relation, err := identity.NewRolePermission(dto.TenantId, dto.RoleId, dto.PermissionId)
	if err != nil {
		return err
	}

	return c.repo.Assign(ctx, relation)
}
