package db

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type ActivateTenantCommandRepository struct {
	db *bun.DB
}

var _ ports.IActivateTenantCommandRepository = (*ActivateTenantCommandRepository)(nil)

func NewActivateTenantCommandRepository(db *bun.DB) *ActivateTenantCommandRepository {
	return &ActivateTenantCommandRepository{db: db}
}

func (r *ActivateTenantCommandRepository) Activate(ctx context.Context, tenantId identity.TenantId) error {
	_, err := r.db.NewUpdate().
		Model((*models.TenantModel)(nil)).
		Set("status = ?", string(identity.TenantStatusActive)).
		Where("id = ?", string(tenantId)).
		Exec(ctx)
	return err
}
