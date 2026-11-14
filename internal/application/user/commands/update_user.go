package commands

import (
	"context"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"
	"kali-auth-context/internal/ports"
	"strings"
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
	if dto.TenantId == "" {
		return identity.ErrTenantRequired
	}
	if strings.TrimSpace(dto.Email) == "" {
		return identity.ErrEmailRequired
	}

	// Contraseña opcional: vacía conserva la actual; si se indica, cumple la política.
	hashedPassword := ""
	if dto.Password != "" {
		if err := policies.ValidatePasswordStrength(dto.Password); err != nil {
			return err
		}
		hashed, err := c.hasher.Hash(dto.Password)
		if err != nil {
			return err
		}
		hashedPassword = hashed
	}

	return c.repo.Update(ctx, &identity.User{
		Id:                   dto.Id,
		TenantId:             dto.TenantId,
		IdentificationNumber: dto.IdentificationNumber,
		Username:             dto.Username,
		Email:                strings.TrimSpace(dto.Email),
		Password:             hashedPassword, // vacío = no se actualiza
	})
}
