package models

import (
	"time"

	"github.com/uptrace/bun"
)

type IdempotencyModel struct {
	bun.BaseModel `bun:"table:idempotency_requests,alias:ir"`
	Key           string    `bun:",pk,type:varchar(255)"`
	Fingerprint   string    `bun:",notnull,type:varchar(255)"`
	StatusCode    int       `bun:",notnull"`
	Body          []byte    `bun:",type:bytea"`
	CreatedAt     time.Time `bun:",notnull,default:current_timestamp"`
}
