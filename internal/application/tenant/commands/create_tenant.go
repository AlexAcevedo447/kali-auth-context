package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type CreateTenantDto struct {
	Name string
}

type CreateTenantCommand struct {
	repo     ports.ICreateTenantCommandRepository
	provider ports.IUUIDProvider
}

func NewCreateTenantCommand(repo ports.ICreateTenantCommandRepository, provider ports.IUUIDProvider) *CreateTenantCommand {
	return &CreateTenantCommand{repo: repo, provider: provider}
}

func (c *CreateTenantCommand) Execute(ctx context.Context, dto *CreateTenantDto) error {
	tenantId := identity.TenantId(c.provider.Generate())
	tenant, err := identity.NewTenant(tenantId, dto.Name)
	if err != nil {
		return err
	}

	return c.repo.Create(ctx, tenant)
}
