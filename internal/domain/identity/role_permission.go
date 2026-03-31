package identity

type RolePermission struct {
	TenantId     TenantId
	RoleId       RoleId
	PermissionId PermissionId
}

func NewRolePermission(tenantId TenantId, roleId RoleId, permissionId PermissionId) (*RolePermission, error) {
	if tenantId == "" {
		return nil, ErrTenantRequired
	}

	if roleId == "" {
		return nil, ErrRoleIdRequired
	}

	if permissionId == "" {
		return nil, ErrPermissionIdRequired
	}

	return &RolePermission{
		TenantId:     tenantId,
		RoleId:       roleId,
		PermissionId: permissionId,
	}, nil
}
