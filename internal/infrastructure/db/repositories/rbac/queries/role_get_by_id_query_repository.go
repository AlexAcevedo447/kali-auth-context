package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type GetRoleByIdQueryRepository struct {
	db *bun.DB
}

var _ ports.IGetRoleByIdQueryRepository = (*GetRoleByIdQueryRepository)(nil)

func NewGetRoleByIdQueryRepository(db *bun.DB) *GetRoleByIdQueryRepository {
	return &GetRoleByIdQueryRepository{db: db}
}

func (r *GetRoleByIdQueryRepository) GetById(ctx context.Context, roleId identity.RoleId) (*identity.Role, error) {
	var model models.RoleModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("id = ?", string(roleId)).
		Scan(ctx); err != nil {
		return nil, err
	}

	return mappers.ToDomainRole(&model), nil
}
