package identity

import (
	"errors"
	"strings"
)

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
	ErrRoleNameTaken    = errors.New("role name already exists")
	ErrRoleAssigned     = errors.New("role is assigned to users")
	ErrProtectedRole    = errors.New("protected administrative role")
)

func NewRole(id RoleId, tenantId TenantId, name, description string) (*Role, error) {
	if id == "" {
		return nil, ErrRoleIdRequired
	}

	if tenantId == "" {
		return nil, ErrTenantRequired
	}

	// Se guarda sin espacios a los lados; solo espacios equivale a vacío.
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrRoleNameRequired
	}

	return &Role{
		Id:          id,
		TenantId:    tenantId,
		Name:        name,
		Description: strings.TrimSpace(description),
	}, nil
}
