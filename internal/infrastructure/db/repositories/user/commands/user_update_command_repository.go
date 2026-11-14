package commands

import (
	"context"
	"database/sql"
	"errors"
	"github.com/uptrace/bun/driver/pgdriver"
	"kali-auth-context/internal/infrastructure/db/models"

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
	// Correo único por empresa (excluyendo al propio usuario).
	taken, err := r.db.NewSelect().
		Model((*models.UserModel)(nil)).
		Where("tenant_id = ? AND email = ? AND id <> ?", user.TenantId, user.Email, user.Id).
		Exists(ctx)
	if err != nil {
		return err
	}
	if taken {
		return identity.ErrEmailAlreadyExists
	}

	model := mappers.ToUserModel(user)
	columns := []string{"identification_number", "username", "email"}
	if model.Password != "" { // contraseña vacía = se conserva la actual
		columns = append(columns, "password")
	}
	res, err := r.db.NewUpdate().
		Model(model).
		Column(columns...).
		Where("tenant_id = ? AND id = ?", model.TenantId, model.Id).
		Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
			return identity.ErrEmailAlreadyExists
		}
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
