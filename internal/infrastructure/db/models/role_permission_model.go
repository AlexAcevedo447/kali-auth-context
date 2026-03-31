package models

import "github.com/uptrace/bun"

type RolePermissionModel struct {
	bun.BaseModel `bun:"table:role_permissions"`

	TenantId     string `bun:"tenant_id,notnull"`
	RoleId       string `bun:"role_id,notnull"`
	PermissionId string `bun:"permission_id,notnull"`
}
