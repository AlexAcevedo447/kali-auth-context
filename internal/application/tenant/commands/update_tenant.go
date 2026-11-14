package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type UpdateTenantDto struct {
	Id     identity.TenantId
	Name   string
	Code   string
	Status identity.TenantStatus
}

type UpdateTenantCommand struct {
	repo ports.IUpdateTenantCommandRepository
}

func NewUpdateTenantCommand(repo ports.IUpdateTenantCommandRepository) *UpdateTenantCommand {
	return &UpdateTenantCommand{repo: repo}
}

func (c *UpdateTenantCommand) Execute(ctx context.Context, dto *UpdateTenantDto) error {
	tenant, err := identity.NewTenant(dto.Id, dto.Name)
	if err != nil {
		return err
	}

	if dto.Code != "" {
		if err := identity.ValidateTenantCode(dto.Code); err != nil {
			return err
		}
		tenant.Code = dto.Code
	}

	if dto.Status == identity.TenantStatusSuspended {
		tenant.Suspend()
	}

	return c.repo.Update(ctx, tenant)
}
