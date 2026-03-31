package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type AssignRoleToUserDto struct {
	TenantId identity.TenantId
	UserId   identity.UserId
	RoleId   identity.RoleId
}

type AssignRoleToUserCommand struct {
	repo ports.IAssignRoleToUserCommandRepository
}

func NewAssignRoleToUserCommand(repo ports.IAssignRoleToUserCommandRepository) *AssignRoleToUserCommand {
	return &AssignRoleToUserCommand{repo: repo}
}

func (c *AssignRoleToUserCommand) Execute(ctx context.Context, dto *AssignRoleToUserDto) error {
	relation, err := identity.NewUserRole(dto.TenantId, dto.UserId, dto.RoleId)
	if err != nil {
		return err
	}

	return c.repo.Assign(ctx, relation)
}
