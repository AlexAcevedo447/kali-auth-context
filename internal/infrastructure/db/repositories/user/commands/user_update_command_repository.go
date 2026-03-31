package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type UpdateUserCommandRepository struct {
	db *bun.DB
}

var _ ports.IUpdateUserCommandRepository = (*UpdateUserCommandRepository)(nil)

func NewUpdateUserCommandRepository(db *bun.DB) *UpdateUserCommandRepository {
	return &UpdateUserCommandRepository{db: db}
}

func (r *UpdateUserCommandRepository) Update(ctx context.Context, user *identity.User) error {
	model := mappers.ToUserModel(user)
	_, err := r.db.NewUpdate().
		Model(model).
		Column("identification_number", "username", "email", "password").
		Where("tenant_id = ? AND id = ?", model.TenantId, model.Id).
		Exec(ctx)
	return err
}
