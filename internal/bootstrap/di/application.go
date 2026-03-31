package bootstrap

import (
	"github.com/google/wire"
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
	"kali-auth-context/internal/ports"
)

var AccessApplicationSet = wire.NewSet(
	authcommands.NewLoginCommand,
	authorizationqueries.NewAuthorizeQuery,
	NewAuthModule,
	NewAuthorizationModule,
)

var UserApplicationSet = wire.NewSet(
	usercommands.NewCreateUserCommand,
	usercommands.NewUpdateUserCommand,
	usercommands.NewDeleteUserCommand,
	userqueries.NewGetUserByIdQuery,
	userqueries.NewGetUserByEmailQuery,
	userqueries.NewListUsersQuery,
	NewUserCommands,
	NewUserQueries,
	NewUserModule,
)

var TenantApplicationSet = wire.NewSet(
	tenantcommands.NewCreateTenantCommand,
	tenantcommands.NewUpdateTenantCommand,
	tenantcommands.NewActivateTenantCommand,
	tenantcommands.NewSuspendTenantCommand,
	tenantqueries.NewGetTenantByIdQuery,
	tenantqueries.NewGetTenantByNameQuery,
	tenantqueries.NewListTenantsQuery,
	NewTenantCommands,
	NewTenantQueries,
	NewTenantModule,
)

var RBACApplicationSet = wire.NewSet(
	rbaccommands.NewCreateRoleCommand,
	rbaccommands.NewUpdateRoleCommand,
	rbaccommands.NewDeleteRoleCommand,
	rbaccommands.NewCreatePermissionCommand,
	rbaccommands.NewUpdatePermissionCommand,
	rbaccommands.NewDeletePermissionCommand,
	rbaccommands.NewAssignRoleToUserCommand,
	rbaccommands.NewRemoveRoleFromUserCommand,
	rbaccommands.NewAssignPermissionToRoleCommand,
	rbaccommands.NewRemovePermissionFromRoleCommand,
	rbacqueries.NewGetRoleByIdQuery,
	rbacqueries.NewListRolesQuery,
	rbacqueries.NewGetPermissionByIdQuery,
	rbacqueries.NewListPermissionsQuery,
	rbacqueries.NewGetUserRolesQuery,
	rbacqueries.NewGetRolePermissionsQuery,
	rbacqueries.NewGetUserEffectivePermissionsQuery,
	NewRBACCommands,
	NewRBACQueries,
	NewRBACModule,
)

func NewSecurityModule(uuidProvider ports.IUUIDProvider, passwordHasher ports.IPasswordHasher) *SecurityModule {
	return &SecurityModule{
		UUIDProvider:   uuidProvider,
		PasswordHasher: passwordHasher,
	}
}

func NewAuthModule(login *authcommands.LoginCommand) *AuthModule {
	return &AuthModule{Login: login}
}

func NewAuthorizationModule(authorize *authorizationqueries.AuthorizeQuery) *AuthorizationModule {
	return &AuthorizationModule{Authorize: authorize}
}

func NewUserCommands(create *usercommands.CreateUserCommand, update *usercommands.UpdateUserCommand, delete *usercommands.DeleteUserCommand) *UserCommands {
	return &UserCommands{Create: create, Update: update, Delete: delete}
}

func NewUserQueries(getById *userqueries.GetUserByIdQuery, getByEmail *userqueries.GetUserByEmailQuery, list *userqueries.ListUsersQuery) *UserQueries {
	return &UserQueries{GetById: getById, GetByEmail: getByEmail, List: list}
}

func NewUserModule(commands *UserCommands, queries *UserQueries) *UserModule {
	return &UserModule{Commands: commands, Queries: queries}
}

func NewTenantCommands(create *tenantcommands.CreateTenantCommand, update *tenantcommands.UpdateTenantCommand, activate *tenantcommands.ActivateTenantCommand, suspend *tenantcommands.SuspendTenantCommand) *TenantCommands {
	return &TenantCommands{Create: create, Update: update, Activate: activate, Suspend: suspend}
}

func NewTenantQueries(getById *tenantqueries.GetTenantByIdQuery, getByName *tenantqueries.GetTenantByNameQuery, list *tenantqueries.ListTenantsQuery) *TenantQueries {
	return &TenantQueries{GetById: getById, GetByName: getByName, List: list}
}

func NewTenantModule(commands *TenantCommands, queries *TenantQueries) *TenantModule {
	return &TenantModule{Commands: commands, Queries: queries}
}

func NewRBACCommands(
	createRole *rbaccommands.CreateRoleCommand,
	updateRole *rbaccommands.UpdateRoleCommand,
	deleteRole *rbaccommands.DeleteRoleCommand,
	createPermission *rbaccommands.CreatePermissionCommand,
	updatePermission *rbaccommands.UpdatePermissionCommand,
	deletePermission *rbaccommands.DeletePermissionCommand,
	assignRoleToUser *rbaccommands.AssignRoleToUserCommand,
	removeRoleFromUser *rbaccommands.RemoveRoleFromUserCommand,
	assignPermissionToRole *rbaccommands.AssignPermissionToRoleCommand,
	removePermissionFromRole *rbaccommands.RemovePermissionFromRoleCommand,
) *RBACCommands {
	return &RBACCommands{
		CreateRole:               createRole,
		UpdateRole:               updateRole,
		DeleteRole:               deleteRole,
		CreatePermission:         createPermission,
		UpdatePermission:         updatePermission,
		DeletePermission:         deletePermission,
		AssignRoleToUser:         assignRoleToUser,
		RemoveRoleFromUser:       removeRoleFromUser,
		AssignPermissionToRole:   assignPermissionToRole,
		RemovePermissionFromRole: removePermissionFromRole,
	}
}

func NewRBACQueries(
	getRoleById *rbacqueries.GetRoleByIdQuery,
	listRoles *rbacqueries.ListRolesQuery,
	getPermissionById *rbacqueries.GetPermissionByIdQuery,
	listPermissions *rbacqueries.ListPermissionsQuery,
	getUserRoles *rbacqueries.GetUserRolesQuery,
	getRolePermissions *rbacqueries.GetRolePermissionsQuery,
	getUserEffectivePermissions *rbacqueries.GetUserEffectivePermissionsQuery,
) *RBACQueries {
	return &RBACQueries{
		GetRoleById:                 getRoleById,
		ListRoles:                   listRoles,
		GetPermissionById:           getPermissionById,
		ListPermissions:             listPermissions,
		GetUserRoles:                getUserRoles,
		GetRolePermissions:          getRolePermissions,
		GetUserEffectivePermissions: getUserEffectivePermissions,
	}
}

func NewRBACModule(commands *RBACCommands, queries *RBACQueries) *RBACModule {
	return &RBACModule{Commands: commands, Queries: queries}
}

func NewCoreContainer(
	cfg *config.Config,
	database *bun.DB,
	security *SecurityModule,
	auth *AuthModule,
	authorization *AuthorizationModule,
	http *HTTPModule,
	user *UserModule,
	tenant *TenantModule,
	rbac *RBACModule,
) *CoreContainer {
	return &CoreContainer{
		Config:        cfg,
		Database:      database,
		Security:      security,
		Auth:          auth,
		Authorization: authorization,
		HTTP:          http,
		User:          user,
		Tenant:        tenant,
		RBAC:          rbac,
	}
}

