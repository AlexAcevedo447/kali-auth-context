package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type DeleteRoleCommand struct {
	repo ports.IDeleteRoleCommandRepository
}

func NewDeleteRoleCommand(repo ports.IDeleteRoleCommandRepository) *DeleteRoleCommand {
	return &DeleteRoleCommand{repo: repo}
}

func (c *DeleteRoleCommand) Execute(ctx context.Context, roleId identity.RoleId) error {
	if roleId == "" {
		return identity.ErrRoleIdRequired
	}
	return c.repo.Delete(ctx, roleId)
}
