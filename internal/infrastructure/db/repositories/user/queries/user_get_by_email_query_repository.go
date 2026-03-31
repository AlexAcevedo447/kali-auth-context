package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type GetUserByEmailQueryRepository struct {
	db *bun.DB
}

var _ ports.IGetUserByEmailQueryRepository = (*GetUserByEmailQueryRepository)(nil)

func NewGetUserByEmailQueryRepository(db *bun.DB) *GetUserByEmailQueryRepository {
	return &GetUserByEmailQueryRepository{db: db}
}

func (r *GetUserByEmailQueryRepository) GetByEmail(tenantId identity.TenantId, email string) (*identity.User, error) {
	var model models.UserModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("tenant_id = ? AND email = ?", tenantId, email).
		Scan(context.Background()); err != nil {
		return nil, err
	}

	return mappers.ToDomainUser(&model)
}
