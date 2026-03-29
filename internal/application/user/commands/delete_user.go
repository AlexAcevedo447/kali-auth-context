package commands

import (
	"context"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type DeleteUserCommand struct {
	repo ports.IDeleteUserCommandRepository
}

func NewDeleteUserCommand(repo ports.IDeleteUserCommandRepository) *DeleteUserCommand {
	return &DeleteUserCommand{repo: repo}
}

func (c *DeleteUserCommand) Execute(ctx context.Context, id *identity.UserId) error {
	return c.repo.Delete(ctx, *id)
}
