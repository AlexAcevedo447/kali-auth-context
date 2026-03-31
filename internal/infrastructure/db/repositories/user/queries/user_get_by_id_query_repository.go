package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type GetUserByIdQueryRepository struct {
	db *bun.DB
}

var _ ports.IGetUserByIdQueryRepository = (*GetUserByIdQueryRepository)(nil)

func NewGetUserByIdQueryRepository(db *bun.DB) *GetUserByIdQueryRepository {
	return &GetUserByIdQueryRepository{db: db}
}

func (r *GetUserByIdQueryRepository) GetById(tenantId identity.TenantId, userId identity.UserId) (*identity.User, error) {
	var model models.UserModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("tenant_id = ? AND id = ?", tenantId, userId).
		Scan(context.Background()); err != nil {
		return nil, err
	}

	return mappers.ToDomainUser(&model)
}
