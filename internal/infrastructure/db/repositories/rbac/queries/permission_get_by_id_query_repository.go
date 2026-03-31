package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type GetPermissionByIdQueryRepository struct {
	db *bun.DB
}

var _ ports.IGetPermissionByIdQueryRepository = (*GetPermissionByIdQueryRepository)(nil)

func NewGetPermissionByIdQueryRepository(db *bun.DB) *GetPermissionByIdQueryRepository {
	return &GetPermissionByIdQueryRepository{db: db}
}

func (r *GetPermissionByIdQueryRepository) GetById(ctx context.Context, permissionId identity.PermissionId) (*identity.Permission, error) {
	var model models.PermissionModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("id = ?", string(permissionId)).
		Scan(ctx); err != nil {
		return nil, err
	}

	return mappers.ToDomainPermission(&model), nil
}
