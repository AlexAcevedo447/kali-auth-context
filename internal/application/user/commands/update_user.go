package commands

import (
	"context"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"
	"kali-auth-context/internal/ports"
)

type UpdateUserDto struct {
	Id                   identity.UserId
	TenantId             identity.TenantId
	IdentificationNumber string
	Username             string
	Email                string
	Password             string
}

type UpdateUserCommand struct {
	repo   ports.IUpdateUserCommandRepository
	hasher ports.IPasswordHasher
}

func NewUpdateUserCommand(repo ports.IUpdateUserCommandRepository, hasher ports.IPasswordHasher) *UpdateUserCommand {
	return &UpdateUserCommand{repo: repo, hasher: hasher}
}

func (c *UpdateUserCommand) Execute(ctx context.Context, dto *UpdateUserDto) error {
	if err := policies.ValidatePasswordStrength(dto.Password); err != nil {
		return err
	}

	hashedPassword, err := c.hasher.Hash(dto.Password)
	if err != nil {
		return err
	}

	user, err := identity.NewUser(
		dto.Id,
		dto.TenantId,
		dto.IdentificationNumber,
		dto.Username,
		dto.Email,
		hashedPassword,
	)
	if err != nil {
		return err
	}

	return c.repo.Update(ctx, user)
}
