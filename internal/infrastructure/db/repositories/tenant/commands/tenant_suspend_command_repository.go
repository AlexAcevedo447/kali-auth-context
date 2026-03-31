package db

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type SuspendTenantCommandRepository struct {
	db *bun.DB
}

var _ ports.ISuspendTenantCommandRepository = (*SuspendTenantCommandRepository)(nil)

func NewSuspendTenantCommandRepository(db *bun.DB) *SuspendTenantCommandRepository {
	return &SuspendTenantCommandRepository{db: db}
}

func (r *SuspendTenantCommandRepository) Suspend(ctx context.Context, tenantId identity.TenantId) error {
	_, err := r.db.NewUpdate().
		Model((*models.TenantModel)(nil)).
		Set("status = ?", string(identity.TenantStatusSuspended)).
		Where("id = ?", string(tenantId)).
		Exec(ctx)
	return err
}
