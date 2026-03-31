package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type UpdateRoleCommandRepository struct {
	db *bun.DB
}

var _ ports.IUpdateRoleCommandRepository = (*UpdateRoleCommandRepository)(nil)

func NewUpdateRoleCommandRepository(db *bun.DB) *UpdateRoleCommandRepository {
	return &UpdateRoleCommandRepository{db: db}
}

func (r *UpdateRoleCommandRepository) Update(ctx context.Context, role *identity.Role) error {
	model := mappers.ToRoleModel(role)
	_, err := r.db.NewUpdate().
		Model(model).
		Column("name", "description").
		Where("id = ?", model.Id).
		Exec(ctx)
	return err
}
