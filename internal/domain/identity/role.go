package identity

import "errors"

type RoleId string

type Role struct {
	Id          RoleId
	TenantId    TenantId
	Name        string
	Description string
}

var (
	ErrRoleIdRequired   = errors.New("role_id is required")
	ErrRoleNameRequired = errors.New("role name is required")
)

func NewRole(id RoleId, tenantId TenantId, name, description string) (*Role, error) {
	if id == "" {
		return nil, ErrRoleIdRequired
	}

	if tenantId == "" {
		return nil, ErrTenantRequired
	}

	if name == "" {
		return nil, ErrRoleNameRequired
	}

	return &Role{
		Id:          id,
		TenantId:    tenantId,
		Name:        name,
		Description: description,
	}, nil
}
