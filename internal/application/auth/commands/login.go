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

// Caso de uso para el inicio de sesión (login) de usuarios.
// Se encarga de validar las credenciales, obtener roles y permisos, y devolver el usuario autenticado.
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

// Ejecuta el proceso de autenticación:
// 1. Valida que los datos no estén vacíos.
// 2. Busca el usuario por email y tenant.
// 3. Si el usuario no existe o la contraseña es incorrecta, retorna error de autenticación.
// 4. Si la autenticación es correcta, obtiene roles y permisos asociados.
// 5. Devuelve el usuario autenticado con sus datos, roles y permisos.
func (c *LoginCommand) Execute(dto *LoginDto) (*AuthenticatedUser, error) {
       if dto == nil || dto.TenantId == "" || dto.Email == "" || dto.Password == "" {
	       // Si faltan datos, retorna error de credenciales inválidas.
	       return nil, identity.ErrInvalidCredentials
       }

       // Busca el usuario por email y tenant.
       user, err := c.userQueryRepo.GetByEmail(dto.TenantId, dto.Email)
       if err != nil {
	       if errors.Is(err, sql.ErrNoRows) {
		       // Si el usuario no existe, simula comparación para evitar timing attacks y retorna error.
		       _ = c.hasher.Compare(c.fakeHash, dto.Password)
		       return nil, identity.ErrInvalidCredentials
	       }
	       return nil, err
       }

       // Compara la contraseña recibida con el hash almacenado.
       if err := c.hasher.Compare(user.Password, dto.Password); err != nil {
	       // Si la contraseña es incorrecta, retorna error de autenticación.
	       return nil, identity.ErrInvalidCredentials
       }

       ctx := context.Background()

       // Obtiene los roles del usuario.
       userRoles, err := c.userRolesRepo.GetByUser(ctx, dto.TenantId, user.Id)
       if err != nil {
	       return nil, err
       }

       roles := make([]*identity.Role, 0, len(userRoles))
       permissionsByID := make(map[identity.PermissionId]*identity.Permission)

       // Por cada rol, obtiene los permisos asociados.
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

       // Prepara la lista de permisos únicos.
       permissions := make([]*identity.Permission, 0, len(permissionsByID))
       for _, perm := range permissionsByID {
	       permissions = append(permissions, perm)
       }

       // Devuelve el usuario autenticado con roles y permisos.
       return &AuthenticatedUser{
	       TenantId:    user.TenantId,
	       UserId:      user.Id,
	       Email:       user.Email,
	       NeedsRehash: c.hasher.NeedsRehash(user.Password),
	       Roles:       roles,
	       Permissions: permissions,
       }, nil
}
