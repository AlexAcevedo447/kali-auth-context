package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type AssignRoleToUserCommandRepository struct {
	db *bun.DB
}

var _ ports.IAssignRoleToUserCommandRepository = (*AssignRoleToUserCommandRepository)(nil)

func NewAssignRoleToUserCommandRepository(db *bun.DB) *AssignRoleToUserCommandRepository {
	return &AssignRoleToUserCommandRepository{db: db}
}

func (r *AssignRoleToUserCommandRepository) Assign(ctx context.Context, relation *identity.UserRole) error {
	model := mappers.ToUserRoleModel(relation)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}
