package idempotency

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type IdempotencyRepository struct {
	db *bun.DB
}

var _ ports.IIdempotencyRepository = (*IdempotencyRepository)(nil)

func NewIdempotencyRepository(db *bun.DB) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

func (r *IdempotencyRepository) Find(ctx context.Context, key string) (*ports.IdempotencyRecord, error) {
	var m models.IdempotencyModel
	err := r.db.NewSelect().
		Model(&m).
		Where("ir.key = ? AND ir.created_at > NOW() - INTERVAL '24 hours'", key).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ports.IdempotencyRecord{
		Key:         m.Key,
		Fingerprint: m.Fingerprint,
		StatusCode:  m.StatusCode,
		Body:        m.Body,
	}, nil
}

func (r *IdempotencyRepository) Save(ctx context.Context, record *ports.IdempotencyRecord) error {
	m := &models.IdempotencyModel{
		Key:         record.Key,
		Fingerprint: record.Fingerprint,
		StatusCode:  record.StatusCode,
		Body:        record.Body,
	}
	_, err := r.db.NewInsert().Model(m).On("CONFLICT (key) DO NOTHING").Exec(ctx)
	return err
}
