package models

import "github.com/uptrace/bun"

type TenantModel struct {
	bun.BaseModel `bun:"table:tenants"`

	Id     string `bun:"id,pk"`
	Name   string `bun:"name,notnull"`
	Status string `bun:"status,notnull"`
}
