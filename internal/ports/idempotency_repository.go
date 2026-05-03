package ports

import "context"

type IdempotencyRecord struct {
	Key         string
	Fingerprint string
	StatusCode  int
	Body        []byte
}

type IIdempotencyRepository interface {
	Find(ctx context.Context, key string) (*IdempotencyRecord, error)
	Save(ctx context.Context, record *IdempotencyRecord) error
}
