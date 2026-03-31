package models

import "github.com/uptrace/bun"

type RoleModel struct {
	bun.BaseModel `bun:"table:roles"`

	Id          string `bun:"id,pk"`
	TenantId    string `bun:"tenant_id,notnull"`
	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`
}
