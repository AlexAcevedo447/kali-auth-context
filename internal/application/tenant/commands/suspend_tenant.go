package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type SuspendTenantCommand struct {
	repo             ports.ISuspendTenantCommandRepository
	platformTenantID identity.TenantId
}

// WithPlatformTenant indica la empresa de plataforma, que no se puede suspender.
func (c *SuspendTenantCommand) WithPlatformTenant(id identity.TenantId) *SuspendTenantCommand {
	c.platformTenantID = id
	return c
}

func NewSuspendTenantCommand(repo ports.ISuspendTenantCommandRepository) *SuspendTenantCommand {
	return &SuspendTenantCommand{repo: repo}
}

func (c *SuspendTenantCommand) Execute(ctx context.Context, tenantId identity.TenantId) error {
	if tenantId == "" {
		return identity.ErrTenantRequired
	}
	if c.platformTenantID != "" && tenantId == c.platformTenantID {
		return identity.ErrPlatformTenantSuspension
	}

	return c.repo.Suspend(ctx, tenantId)
}
