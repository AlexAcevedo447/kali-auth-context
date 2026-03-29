package commands

import (
	"context"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type CreateUserDto struct {
	IdentificationNumber string
	Username             string
	Email                string
	Password             string
}

type CreateUserCommand struct {
	repo     ports.ICreateUserCommandRepository
	provider ports.IUUIDProvider
	hasher   ports.IPasswordHasher
}


func NewCreateUserCommand(repo ports.ICreateUserCommandRepository, provider ports.IUUIDProvider, hasher ports.IPasswordHasher) *CreateUserCommand {
	return &CreateUserCommand{repo: repo, provider: provider, hasher: hasher}
}

func (c *CreateUserCommand) Execute(ctx context.Context, user *CreateUserDto) error {
	id := identity.UserId(c.provider.Generate())

	hashedPassword, err := c.hasher.Hash(user.Password)

	if err != nil {
		return err
	}

	domainUser, err := identity.NewUser(
		id,
		user.IdentificationNumber,
		user.Username,
		user.Email,
		hashedPassword,
	)

	if err != nil {
		return err
	}

	return c.repo.Create(ctx, domainUser)
}
