package commands

import (
	"context"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type DeleteUserCommand struct {
	repo ports.IDeleteUserCommandRepository
}

type DeleteUserDto struct {
	TenantId identity.TenantId
	UserId   identity.UserId
}

func NewDeleteUserCommand(repo ports.IDeleteUserCommandRepository) *DeleteUserCommand {
	return &DeleteUserCommand{repo: repo}
}

func (c *DeleteUserCommand) Execute(ctx context.Context, dto *DeleteUserDto) error {
	if dto == nil {
		return identity.ErrUserIdRequired
	}
	if dto.TenantId == "" {
		return identity.ErrTenantRequired
	}
	if dto.UserId == "" {
		return identity.ErrUserIdRequired
	}

	return c.repo.Delete(ctx, dto.TenantId, dto.UserId)
}
