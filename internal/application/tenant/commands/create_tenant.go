package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type CreateTenantDto struct {
	Name string
	// Code es el código legible de la empresa; si no llega se deriva del nombre.
	Code string
}

type CreateTenantCommand struct {
	repo     ports.ICreateTenantCommandRepository
	provider ports.IUUIDProvider
}

func NewCreateTenantCommand(repo ports.ICreateTenantCommandRepository, provider ports.IUUIDProvider) *CreateTenantCommand {
	return &CreateTenantCommand{repo: repo, provider: provider}
}

func (c *CreateTenantCommand) Execute(ctx context.Context, dto *CreateTenantDto) (*identity.Tenant, error) {
	tenantId := identity.TenantId(c.provider.Generate())
	tenant, err := identity.NewTenant(tenantId, dto.Name)
	if err != nil {
		return nil, err
	}

	tenant.Code = dto.Code
	if tenant.Code == "" {
		tenant.Code = identity.SlugifyTenantCode(dto.Name)
	}
	if err := identity.ValidateTenantCode(tenant.Code); err != nil {
		return nil, err
	}

	if err := c.repo.Create(ctx, tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}
