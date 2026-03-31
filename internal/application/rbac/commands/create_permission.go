package commands

import (
	"context"

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

func (c *CreatePermissionCommand) Execute(ctx context.Context, dto *CreatePermissionDto) error {
	permissionId := identity.PermissionId(c.provider.Generate())
	permission, err := identity.NewPermission(permissionId, dto.TenantId, dto.Resource, dto.Action)
	if err != nil {
		return err
	}

	return c.repo.Create(ctx, permission)
}
