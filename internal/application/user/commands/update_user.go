package commands

import (
	"context"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type UpdateUserDto struct {
	IdentificationNumber string
	Username             string
	Email                string
	Password             string
}

type UpdateUserCommand struct {
	repo ports.IUpdateUserCommandRepository
}

func NewUpdateUserCommand(repo ports.IUpdateUserCommandRepository) *UpdateUserCommand {
	return &UpdateUserCommand{repo: repo}
}

func (c *UpdateUserCommand) Execute(ctx context.Context, dto *UpdateUserDto) error {
	user := &identity.User{
		IdentificationNumber: dto.IdentificationNumber,
		Username:             dto.Username,
		Email:                dto.Email,
		Password:             dto.Password,
	}

	return c.repo.Update(ctx, user)
}
