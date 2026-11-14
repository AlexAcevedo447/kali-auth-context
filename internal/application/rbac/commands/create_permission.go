package commands

import (
	"context"
	"kali-auth-context/internal/domain/policies"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type CreatePermissionDto struct {
	TenantId identity.TenantId
	Resource string
	Action   string
}

type CreatePermissionCommand struct {
	repo     ports.ICreatePermissionCommandRepository
	provider ports.IUUIDProvider
}

func NewCreatePermissionCommand(repo ports.ICreatePermissionCommandRepository, provider ports.IUUIDProvider) *CreatePermissionCommand {
	return &CreatePermissionCommand{repo: repo, provider: provider}
}

func (c *CreatePermissionCommand) Execute(ctx context.Context, dto *CreatePermissionDto) (*identity.Permission, error) {
	// Normaliza (minúsculas, sin espacios) y valida contra el catálogo.
	resource, action, err := policies.ValidatePermission(dto.Resource, dto.Action)
	if err != nil {
		return nil, err
	}
	permissionId := identity.PermissionId(c.provider.Generate())
	permission, err := identity.NewPermission(permissionId, dto.TenantId, resource, action)
	if err != nil {
		return nil, err
	}

	if err := c.repo.Create(ctx, permission); err != nil {
		return nil, err
	}
	return permission, nil
}
