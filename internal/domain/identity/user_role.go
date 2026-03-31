package identity

import "errors"

type UserRole struct {
	TenantId TenantId
	UserId   UserId
	RoleId   RoleId
}

var (
ErrUserIdRequired = errors.New("user_id is required")
)

func NewUserRole(tenantId TenantId, userId UserId, roleId RoleId) (*UserRole, error) {
	if tenantId == "" {
		return nil, ErrTenantRequired
	}

	if userId == "" {
		return nil, ErrUserIdRequired
	}

	if roleId == "" {
		return nil, ErrRoleIdRequired
	}

	return &UserRole{
		TenantId: tenantId,
		UserId:   userId,
		RoleId:   roleId,
	}, nil
}
