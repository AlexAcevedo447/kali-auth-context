package models

import "github.com/uptrace/bun"

type TenantModel struct {
	bun.BaseModel `bun:"table:tenants"`

	Id     string `bun:"id,pk"`
	Name   string `bun:"name,notnull"`
	Code   string `bun:"code,nullzero"`
	Status string `bun:"status,notnull"`
}
