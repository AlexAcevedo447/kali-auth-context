package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type ListUsersQueryRepository struct {
	db *bun.DB
}

var _ ports.IListUsersQueryRepository = (*ListUsersQueryRepository)(nil)

func NewListUsersQueryRepository(db *bun.DB) *ListUsersQueryRepository {
	return &ListUsersQueryRepository{db: db}
}

func (r *ListUsersQueryRepository) List(tenantId identity.TenantId) ([]*identity.User, error) {
	var modelsList []models.UserModel
	if err := r.db.NewSelect().
		Model(&modelsList).
		Where("tenant_id = ?", tenantId).
		Scan(context.Background()); err != nil {
		return nil, err
	}

	result := make([]*identity.User, 0, len(modelsList))
	for i := range modelsList {
		u, err := mappers.ToDomainUser(&modelsList[i])
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	return result, nil
}
