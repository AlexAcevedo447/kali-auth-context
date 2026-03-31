package models

import "github.com/uptrace/bun"

type PermissionModel struct {
	bun.BaseModel `bun:"table:permissions"`

	Id       string `bun:"id,pk"`
	TenantId string `bun:"tenant_id,notnull"`
	Resource string `bun:"resource,notnull"`
	Action   string `bun:"action,notnull"`
}
