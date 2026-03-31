package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type DeletePermissionCommand struct {
	repo ports.IDeletePermissionCommandRepository
}

func NewDeletePermissionCommand(repo ports.IDeletePermissionCommandRepository) *DeletePermissionCommand {
	return &DeletePermissionCommand{repo: repo}
}

func (c *DeletePermissionCommand) Execute(ctx context.Context, permissionId identity.PermissionId) error {
	if permissionId == "" {
		return identity.ErrPermissionIdRequired
	}
	return c.repo.Delete(ctx, permissionId)
}
