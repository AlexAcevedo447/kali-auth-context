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
	repo             ports.IUpdateTenantCommandRepository
	platformTenantID identity.TenantId
}

// WithPlatformTenant indica la empresa de plataforma, que no se puede suspender
// tampoco por PUT.
func (c *UpdateTenantCommand) WithPlatformTenant(id identity.TenantId) *UpdateTenantCommand {
	c.platformTenantID = id
	return c
}

func NewUpdateTenantCommand(repo ports.IUpdateTenantCommandRepository) *UpdateTenantCommand {
	return &UpdateTenantCommand{repo: repo}
}

func (c *UpdateTenantCommand) Execute(ctx context.Context, dto *UpdateTenantDto) error {
	// Un estado vacío conserva el comportamiento previo (queda ACTIVE); cualquier
	// otro valor fuera del dominio es un error de validación.
	switch dto.Status {
	case "", identity.TenantStatusActive, identity.TenantStatusSuspended:
	default:
		return identity.ErrTenantStatusInvalid
	}

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
		if c.platformTenantID != "" && dto.Id == c.platformTenantID {
			return identity.ErrPlatformTenantSuspension
		}
		tenant.Suspend()
	}

	return c.repo.Update(ctx, tenant)
}
