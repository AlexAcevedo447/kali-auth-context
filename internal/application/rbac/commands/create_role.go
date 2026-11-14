package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type CreateRoleDto struct {
	TenantId    identity.TenantId
	Name        string
	Description string
}

type CreateRoleCommand struct {
	repo     ports.ICreateRoleCommandRepository
	provider ports.IUUIDProvider
}

func NewCreateRoleCommand(repo ports.ICreateRoleCommandRepository, provider ports.IUUIDProvider) *CreateRoleCommand {
	return &CreateRoleCommand{repo: repo, provider: provider}
}

func (c *CreateRoleCommand) Execute(ctx context.Context, dto *CreateRoleDto) (*identity.Role, error) {
	roleId := identity.RoleId(c.provider.Generate())
	role, err := identity.NewRole(roleId, dto.TenantId, dto.Name, dto.Description)
	if err != nil {
		return nil, err
	}

	if err := c.repo.Create(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}
