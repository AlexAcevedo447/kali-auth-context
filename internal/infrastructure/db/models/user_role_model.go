package models

import "github.com/uptrace/bun"

type UserRoleModel struct {
	bun.BaseModel `bun:"table:user_roles"`

	TenantId string `bun:"tenant_id,notnull"`
	UserId   string `bun:"user_id,notnull"`
	RoleId   string `bun:"role_id,notnull"`
}
