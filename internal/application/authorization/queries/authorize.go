package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"
	"kali-auth-context/internal/ports"
)

type AuthorizeDto struct {
	TenantId identity.TenantId
	UserId   identity.UserId
	Resource string
	Action   string
}

type AuthorizationDecision struct {
	Allowed bool
	Reason  string
}

type AuthorizeQuery struct {
	tenantRepo          ports.IGetTenantByIdQueryRepository
	userRepo            ports.IGetUserByIdQueryRepository
	userRolesRepo       ports.IGetUserRolesQueryRepository
	rolePermissionsRepo ports.IGetRolePermissionsQueryRepository
	permissionRepo      ports.IGetPermissionByIdQueryRepository
}

func NewAuthorizeQuery(
	tenantRepo ports.IGetTenantByIdQueryRepository,
	userRepo ports.IGetUserByIdQueryRepository,
	userRolesRepo ports.IGetUserRolesQueryRepository,
	rolePermissionsRepo ports.IGetRolePermissionsQueryRepository,
	permissionRepo ports.IGetPermissionByIdQueryRepository,
) *AuthorizeQuery {
	return &AuthorizeQuery{
		tenantRepo:          tenantRepo,
		userRepo:            userRepo,
		userRolesRepo:       userRolesRepo,
		rolePermissionsRepo: rolePermissionsRepo,
		permissionRepo:      permissionRepo,
	}
}

func (q *AuthorizeQuery) Handle(ctx context.Context, dto *AuthorizeDto) (*AuthorizationDecision, error) {
	if dto == nil {
		return nil, identity.ErrAuthorizationRequestRequired
	}
	if dto.TenantId == "" {
		return nil, identity.ErrTenantRequired
	}
	if dto.UserId == "" {
		return nil, identity.ErrUserIdRequired
	}
	if dto.Resource == "" {
		return nil, identity.ErrAuthorizationResourceRequired
	}
	if dto.Action == "" {
		return nil, identity.ErrAuthorizationActionRequired
	}

	tenant, err := q.tenantRepo.GetById(ctx, dto.TenantId)
	if err != nil {
		return nil, err
	}
	if !tenant.IsActive() {
		return &AuthorizationDecision{Allowed: false, Reason: "tenant is suspended"}, nil
	}

	if _, err := q.userRepo.GetById(dto.TenantId, dto.UserId); err != nil {
		return &AuthorizationDecision{Allowed: false, Reason: "user is not assigned to tenant"}, nil
	}

	permissions, err := q.getUserEffectivePermissions(ctx, dto.TenantId, dto.UserId)
	if err != nil {
		return nil, err
	}

	if !policies.HasPermission(dto.TenantId, permissions, dto.Resource, dto.Action) {
		return &AuthorizationDecision{Allowed: false, Reason: "permission denied"}, nil
	}

	return &AuthorizationDecision{Allowed: true}, nil
}

func (q *AuthorizeQuery) getUserEffectivePermissions(ctx context.Context, tenantId identity.TenantId, userId identity.UserId) ([]*identity.Permission, error) {
	userRoles, err := q.userRolesRepo.GetByUser(ctx, tenantId, userId)
	if err != nil {
		return nil, err
	}

	permissionById := make(map[identity.PermissionId]*identity.Permission)
	for _, userRole := range userRoles {
		rolePermissions, err := q.rolePermissionsRepo.GetByRole(ctx, tenantId, userRole.RoleId)
		if err != nil {
			return nil, err
		}

		for _, rolePermission := range rolePermissions {
			if _, exists := permissionById[rolePermission.PermissionId]; exists {
				continue
			}

			permission, err := q.permissionRepo.GetById(ctx, rolePermission.PermissionId)
			if err != nil {
				return nil, err
			}

			permissionById[rolePermission.PermissionId] = permission
		}
	}

	result := make([]*identity.Permission, 0, len(permissionById))
	for _, permission := range permissionById {
		result = append(result, permission)
	}

	return result, nil
}
