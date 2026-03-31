package bootstrap

import (
	"github.com/uptrace/bun"

	authcommands "kali-auth-context/internal/application/auth/commands"
	authorizationqueries "kali-auth-context/internal/application/authorization/queries"
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	rbacqueries "kali-auth-context/internal/application/rbac/queries"
	tenantcommands "kali-auth-context/internal/application/tenant/commands"
	tenantqueries "kali-auth-context/internal/application/tenant/queries"
	usercommands "kali-auth-context/internal/application/user/commands"
	userqueries "kali-auth-context/internal/application/user/queries"
	"kali-auth-context/internal/infrastructure/config"
	"kali-auth-context/internal/infrastructure/http/routes"
	"kali-auth-context/internal/ports"
)

type CoreContainer struct {
	Config        *config.Config
	Database      *bun.DB
	Security      *SecurityModule
	Auth          *AuthModule
	Authorization *AuthorizationModule
	HTTP          *HTTPModule
	User          *UserModule
	Tenant        *TenantModule
	RBAC          *RBACModule
}

func (c *CoreContainer) Close() error {
	if c == nil || c.Database == nil {
		return nil
	}

	return c.Database.Close()
}

type SecurityModule struct {
	UUIDProvider   ports.IUUIDProvider
	PasswordHasher ports.IPasswordHasher
}

type AuthModule struct {
	Login *authcommands.LoginCommand
}

type AuthorizationModule struct {
	Authorize *authorizationqueries.AuthorizeQuery
}

type HTTPModule struct {
	Router *routes.Router
}

type UserModule struct {
	Commands *UserCommands
	Queries  *UserQueries
}

type UserCommands struct {
	Create *usercommands.CreateUserCommand
	Update *usercommands.UpdateUserCommand
	Delete *usercommands.DeleteUserCommand
}

type UserQueries struct {
	GetById    *userqueries.GetUserByIdQuery
	GetByEmail *userqueries.GetUserByEmailQuery
	List       *userqueries.ListUsersQuery
}

type TenantModule struct {
	Commands *TenantCommands
	Queries  *TenantQueries
}

type TenantCommands struct {
	Create   *tenantcommands.CreateTenantCommand
	Update   *tenantcommands.UpdateTenantCommand
	Activate *tenantcommands.ActivateTenantCommand
	Suspend  *tenantcommands.SuspendTenantCommand
}

type TenantQueries struct {
	GetById   *tenantqueries.GetTenantByIdQuery
	GetByName *tenantqueries.GetTenantByNameQuery
	List      *tenantqueries.ListTenantsQuery
}

type RBACModule struct {
	Commands *RBACCommands
	Queries  *RBACQueries
}

type RBACCommands struct {
	CreateRole               *rbaccommands.CreateRoleCommand
	UpdateRole               *rbaccommands.UpdateRoleCommand
	DeleteRole               *rbaccommands.DeleteRoleCommand
	CreatePermission         *rbaccommands.CreatePermissionCommand
	UpdatePermission         *rbaccommands.UpdatePermissionCommand
	DeletePermission         *rbaccommands.DeletePermissionCommand
	AssignRoleToUser         *rbaccommands.AssignRoleToUserCommand
	RemoveRoleFromUser       *rbaccommands.RemoveRoleFromUserCommand
	AssignPermissionToRole   *rbaccommands.AssignPermissionToRoleCommand
	RemovePermissionFromRole *rbaccommands.RemovePermissionFromRoleCommand
}

type RBACQueries struct {
	GetRoleById               *rbacqueries.GetRoleByIdQuery
	ListRoles                 *rbacqueries.ListRolesQuery
	GetPermissionById         *rbacqueries.GetPermissionByIdQuery
	ListPermissions           *rbacqueries.ListPermissionsQuery
	GetUserRoles              *rbacqueries.GetUserRolesQuery
	GetRolePermissions        *rbacqueries.GetRolePermissionsQuery
	GetUserEffectivePermissions *rbacqueries.GetUserEffectivePermissionsQuery
}

