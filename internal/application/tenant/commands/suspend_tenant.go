package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type SuspendTenantCommand struct {
	repo ports.ISuspendTenantCommandRepository
}

func NewSuspendTenantCommand(repo ports.ISuspendTenantCommandRepository) *SuspendTenantCommand {
	return &SuspendTenantCommand{repo: repo}
}

func (c *SuspendTenantCommand) Execute(ctx context.Context, tenantId identity.TenantId) error {
	if tenantId == "" {
		return identity.ErrTenantRequired
	}

	return c.repo.Suspend(ctx, tenantId)
}
