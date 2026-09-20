package commands

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

// postgresUniqueViolation es el código SQLSTATE de Postgres para
// "unique_violation" (23505). Se usa para traducir la violación del índice
// único idx_users_tenant_id_email_unique a un error de dominio.
const postgresUniqueViolation = "23505"

type CreateUserCommandRepository struct {
	db *bun.DB
}

var _ ports.ICreateUserCommandRepository = (*CreateUserCommandRepository)(nil)

func NewCreateUserCommandRepository(db *bun.DB) *CreateUserCommandRepository {
	return &CreateUserCommandRepository{db: db}
}

func (r *CreateUserCommandRepository) Create(ctx context.Context, user *identity.User) error {
	// Prevalidación: entrega un error de dominio claro en el caso común,
	// sin depender de ella como única protección (ver el índice único
	// idx_users_tenant_id_email_unique, capturado más abajo).
	exists, err := r.db.NewSelect().
		Model((*models.UserModel)(nil)).
		Where("tenant_id = ? AND email = ?", user.TenantId, user.Email).
		Exists(ctx)
	if err != nil {
		return err
	}
	if exists {
		return identity.ErrEmailAlreadyExists
	}

	model := mappers.ToUserModel(user)
	if _, err := r.db.NewInsert().Model(model).Exec(ctx); err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == postgresUniqueViolation {
			return identity.ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}
