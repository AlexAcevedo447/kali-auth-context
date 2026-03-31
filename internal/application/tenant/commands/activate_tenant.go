package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type ActivateTenantCommand struct {
	repo ports.IActivateTenantCommandRepository
}

func NewActivateTenantCommand(repo ports.IActivateTenantCommandRepository) *ActivateTenantCommand {
	return &ActivateTenantCommand{repo: repo}
}

func (c *ActivateTenantCommand) Execute(ctx context.Context, tenantId identity.TenantId) error {
	if tenantId == "" {
		return identity.ErrTenantRequired
	}

	return c.repo.Activate(ctx, tenantId)
}
