package commands

import (
	"database/sql"
	"errors"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type LoginDto struct {
	TenantId identity.TenantId
	Email    string
	Password string
}

type AuthenticatedUser struct {
	TenantId    identity.TenantId
	UserId      identity.UserId
	Email       string
	NeedsRehash bool
}

type LoginCommand struct {
	userQueryRepo ports.IGetUserByEmailQueryRepository
	hasher        ports.IPasswordHasher
	fakeHash      string
}

func NewLoginCommand(userQueryRepo ports.IGetUserByEmailQueryRepository, hasher ports.IPasswordHasher) *LoginCommand {
	fakeHash, err := hasher.Hash("fake-password-for-timing-equalization")
	if err != nil {
		fakeHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoOe6j6M1Kf3r5L6N8A1sC3yRyvE4s46aW"
	}

	return &LoginCommand{
		userQueryRepo: userQueryRepo,
		hasher:        hasher,
		fakeHash:      fakeHash,
	}
}

func (c *LoginCommand) Execute(dto *LoginDto) (*AuthenticatedUser, error) {
	if dto == nil || dto.TenantId == "" || dto.Email == "" || dto.Password == "" {
		return nil, identity.ErrInvalidCredentials
	}

	user, err := c.userQueryRepo.GetByEmail(dto.TenantId, dto.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.hasher.Compare(c.fakeHash, dto.Password)
			return nil, identity.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := c.hasher.Compare(user.Password, dto.Password); err != nil {
		return nil, identity.ErrInvalidCredentials
	}

	return &AuthenticatedUser{
		TenantId:    user.TenantId,
		UserId:      user.Id,
		Email:       user.Email,
		NeedsRehash: c.hasher.NeedsRehash(user.Password),
	}, nil
}
