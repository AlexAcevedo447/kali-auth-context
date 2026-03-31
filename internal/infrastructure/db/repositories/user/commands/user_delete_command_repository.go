package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type DeleteUserCommandRepository struct {
	db *bun.DB
}

var _ ports.IDeleteUserCommandRepository = (*DeleteUserCommandRepository)(nil)

func NewDeleteUserCommandRepository(db *bun.DB) *DeleteUserCommandRepository {
	return &DeleteUserCommandRepository{db: db}
}

func (r *DeleteUserCommandRepository) Delete(ctx context.Context, tenantId identity.TenantId, userId identity.UserId) error {
	_, err := r.db.NewDelete().
		Model((*models.UserModel)(nil)).
		Where("tenant_id = ? AND id = ?", tenantId, userId).
		Exec(ctx)
	return err
}
