package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type RemoveRoleFromUserDto struct {
	TenantId identity.TenantId
	UserId   identity.UserId
	RoleId   identity.RoleId
}

type RemoveRoleFromUserCommand struct {
	repo ports.IRemoveRoleFromUserCommandRepository
}

func NewRemoveRoleFromUserCommand(repo ports.IRemoveRoleFromUserCommandRepository) *RemoveRoleFromUserCommand {
	return &RemoveRoleFromUserCommand{repo: repo}
}

func (c *RemoveRoleFromUserCommand) Execute(ctx context.Context, dto *RemoveRoleFromUserDto) error {
	if dto.TenantId == "" {
		return identity.ErrTenantRequired
	}
	if dto.UserId == "" {
		return identity.ErrUserIdRequired
	}
	if dto.RoleId == "" {
		return identity.ErrRoleIdRequired
	}

	return c.repo.Remove(ctx, dto.TenantId, dto.UserId, dto.RoleId)
}
