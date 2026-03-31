package identity

import "errors"

type PermissionId string

type Permission struct {
	Id       PermissionId
	TenantId TenantId
	Resource string
	Action   string
}

var (
	ErrPermissionIdRequired = errors.New("permission_id is required")
	ErrPermissionResource   = errors.New("permission resource is required")
	ErrPermissionAction     = errors.New("permission action is required")
)

func NewPermission(id PermissionId, tenantId TenantId, resource, action string) (*Permission, error) {
	if id == "" {
		return nil, ErrPermissionIdRequired
	}

	if tenantId == "" {
		return nil, ErrTenantRequired
	}

	if resource == "" {
		return nil, ErrPermissionResource
	}

	if action == "" {
		return nil, ErrPermissionAction
	}

	return &Permission{
		Id:       id,
		TenantId: tenantId,
		Resource: resource,
		Action:   action,
	}, nil
}
