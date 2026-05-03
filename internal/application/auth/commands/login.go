package commands

import (
	"context"
	"database/sql"
	"errors"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type LoginDto struct {
	TenantId identity.TenantId
	Email    string
	Password string
}

type AuthenticatedUser struct {
	TenantId    identity.TenantId
	UserId      identity.UserId
	Email       string
	NeedsRehash bool
	Roles       []*identity.Role
	Permissions []*identity.Permission
}

type LoginCommand struct {
	userQueryRepo           ports.IGetUserByEmailQueryRepository
	userRolesRepo           ports.IGetUserRolesQueryRepository
	roleByIdRepo            ports.IGetRoleByIdQueryRepository
	rolePermissionsRepo     ports.IGetRolePermissionsQueryRepository
	permissionQueryRepo     ports.IGetPermissionByIdQueryRepository
	hasher                  ports.IPasswordHasher
	fakeHash                string
}

func NewLoginCommand(
	userQueryRepo ports.IGetUserByEmailQueryRepository,
	userRolesRepo ports.IGetUserRolesQueryRepository,
	roleByIdRepo ports.IGetRoleByIdQueryRepository,
	rolePermissionsRepo ports.IGetRolePermissionsQueryRepository,
	permissionQueryRepo ports.IGetPermissionByIdQueryRepository,
	hasher ports.IPasswordHasher,
) *LoginCommand {
	fakeHash, err := hasher.Hash("fake-password-for-timing-equalization")
	if err != nil {
		fakeHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoOe6j6M1Kf3r5L6N8A1sC3yRyvE4s46aW"
	}

	return &LoginCommand{
		userQueryRepo:       userQueryRepo,
		userRolesRepo:       userRolesRepo,
		roleByIdRepo:        roleByIdRepo,
		rolePermissionsRepo: rolePermissionsRepo,
		permissionQueryRepo: permissionQueryRepo,
		hasher:              hasher,
		fakeHash:            fakeHash,
	}
}

func (c *LoginCommand) Execute(dto *LoginDto) (*AuthenticatedUser, error) {
	if dto == nil || dto.TenantId == "" || dto.Email == "" || dto.Password == "" {
		return nil, identity.ErrInvalidCredentials
	}

	user, err := c.userQueryRepo.GetByEmail(dto.TenantId, dto.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = c.hasher.Compare(c.fakeHash, dto.Password)
			return nil, identity.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := c.hasher.Compare(user.Password, dto.Password); err != nil {
		return nil, identity.ErrInvalidCredentials
	}

	ctx := context.Background()

	userRoles, err := c.userRolesRepo.GetByUser(ctx, dto.TenantId, user.Id)
	if err != nil {
		return nil, err
	}

	roles := make([]*identity.Role, 0, len(userRoles))
	permissionsByID := make(map[identity.PermissionId]*identity.Permission)

	for _, userRole := range userRoles {
		role, err := c.roleByIdRepo.GetById(ctx, userRole.RoleId)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)

		rolePerms, err := c.rolePermissionsRepo.GetByRole(ctx, dto.TenantId, userRole.RoleId)
		if err != nil {
			return nil, err
		}

		for _, rolePerm := range rolePerms {
			if _, exists := permissionsByID[rolePerm.PermissionId]; exists {
				continue
			}

			perm, err := c.permissionQueryRepo.GetById(ctx, rolePerm.PermissionId)
			if err != nil {
				return nil, err
			}

			permissionsByID[rolePerm.PermissionId] = perm
		}
	}

	permissions := make([]*identity.Permission, 0, len(permissionsByID))
	for _, perm := range permissionsByID {
		permissions = append(permissions, perm)
	}

	return &AuthenticatedUser{
		TenantId:    user.TenantId,
		UserId:      user.Id,
		Email:       user.Email,
		NeedsRehash: c.hasher.NeedsRehash(user.Password),
		Roles:       roles,
		Permissions: permissions,
	}, nil
}
