package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

// LoginDto son los datos que escribe una persona: correo y contraseña. El código
// de empresa solo se envía cuando el servidor lo pide (cuentas en varias empresas).
// El tenant nunca llega del cliente: se resuelve a partir de la cuenta.
type LoginDto struct {
	Email      string
	Password   string
	TenantCode string
}

type AuthenticatedUser struct {
	TenantId    identity.TenantId
	TenantName  string
	TenantCode  string
	UserId      identity.UserId
	Email       string
	NeedsRehash bool
	Roles       []*identity.Role
	Permissions []*identity.Permission
}

type LoginCommand struct {
	candidatesRepo      ports.IFindLoginCandidatesQueryRepository
	userRolesRepo       ports.IGetUserRolesQueryRepository
	roleByIdRepo        ports.IGetRoleByIdQueryRepository
	rolePermissionsRepo ports.IGetRolePermissionsQueryRepository
	permissionQueryRepo ports.IGetPermissionByIdQueryRepository
	hasher              ports.IPasswordHasher
	fakeHash            string
}

func NewLoginCommand(
	candidatesRepo ports.IFindLoginCandidatesQueryRepository,
	userRolesRepo ports.IGetUserRolesQueryRepository,
	roleByIdRepo ports.IGetRoleByIdQueryRepository,
	rolePermissionsRepo ports.IGetRolePermissionsQueryRepository,
	permissionQueryRepo ports.IGetPermissionByIdQueryRepository,
	hasher ports.IPasswordHasher,
) *LoginCommand {
	// Hash ficticio para igualar tiempos cuando el correo no existe.
	fakeHash, err := hasher.Hash("fake-password-for-timing-equalization")
	if err != nil {
		fakeHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoOe6j6M1Kf3r5L6N8A1sC3yRyvE4s46aW"
	}

	return &LoginCommand{
		candidatesRepo:      candidatesRepo,
		userRolesRepo:       userRolesRepo,
		roleByIdRepo:        roleByIdRepo,
		rolePermissionsRepo: rolePermissionsRepo,
		permissionQueryRepo: permissionQueryRepo,
		hasher:              hasher,
		fakeHash:            fakeHash,
	}
}

func (c *LoginCommand) Execute(dto *LoginDto) (*AuthenticatedUser, error) {
	if dto == nil || dto.Email == "" || dto.Password == "" {
		return nil, identity.ErrInvalidCredentials
	}

	ctx := context.Background()

	candidates, err := c.candidatesRepo.FindByEmail(ctx, dto.Email, dto.TenantCode)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		_ = c.hasher.Compare(c.fakeHash, dto.Password)
		return nil, identity.ErrInvalidCredentials
	}

	// El tenant se resuelve por la cuenta cuya contraseña coincide. Si coinciden
	// cuentas de varias empresas, se pide el código de empresa (nunca el UUID).
	matches := make([]ports.LoginCandidate, 0, 1)
	for _, candidate := range candidates {
		if err := c.hasher.Compare(candidate.User.Password, dto.Password); err == nil {
			matches = append(matches, candidate)
		}
	}
	switch {
	case len(matches) == 0:
		return nil, identity.ErrInvalidCredentials
	case len(matches) > 1:
		return nil, identity.ErrTenantCodeRequired
	}

	match := matches[0]
	if match.TenantStatus != identity.TenantStatusActive {
		return nil, identity.ErrTenantSuspended
	}
	user := match.User

	userRoles, err := c.userRolesRepo.GetByUser(ctx, user.TenantId, user.Id)
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
		// Defensa en profundidad: un rol de otra empresa nunca concede acceso.
		if role.TenantId != user.TenantId {
			continue
		}
		roles = append(roles, role)

		rolePerms, err := c.rolePermissionsRepo.GetByRole(ctx, user.TenantId, userRole.RoleId)
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
			if perm.TenantId != user.TenantId {
				continue
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
		TenantName:  match.TenantName,
		TenantCode:  match.TenantCode,
		UserId:      user.Id,
		Email:       user.Email,
		NeedsRehash: c.hasher.NeedsRehash(user.Password),
		Roles:       roles,
		Permissions: permissions,
	}, nil
}
