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
	ErrPermissionTaken      = errors.New("permission already exists")
	ErrPermissionAssigned   = errors.New("permission is assigned to roles")
	ErrProtectedPermission  = errors.New("protected administrative permission")
	ErrUnknownResource      = errors.New("unknown permission resource")
	ErrUnknownAction        = errors.New("unknown permission action")
	// Autorización RBAC (403) y protección del administrador (409).
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrPrivilegeEscalation     = errors.New("cannot grant privileges you do not have")
	ErrLastAdministrator       = errors.New("operation would leave the tenant without an administrator")
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
