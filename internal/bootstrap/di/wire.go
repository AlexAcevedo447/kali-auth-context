package bootstrap

import (
	"context"
	"time"

	"github.com/google/wire"

	authcommands "kali-auth-context/internal/application/auth/commands"
	authorizationqueries "kali-auth-context/internal/application/authorization/queries"
	rbaccommands "kali-auth-context/internal/application/rbac/commands"
	rbacqueries "kali-auth-context/internal/application/rbac/queries"
	tenantcommands "kali-auth-context/internal/application/tenant/commands"
	tenantqueries "kali-auth-context/internal/application/tenant/queries"
	usercommands "kali-auth-context/internal/application/user/commands"
	userqueries "kali-auth-context/internal/application/user/queries"
	"kali-auth-context/internal/infrastructure/config"
	"kali-auth-context/internal/infrastructure/db"
	idempotencyrepo "kali-auth-context/internal/infrastructure/db/repositories/idempotency"
	rbaccmdrepo "kali-auth-context/internal/infrastructure/db/repositories/rbac/commands"
	rbacqryrepo "kali-auth-context/internal/infrastructure/db/repositories/rbac/queries"
	tenantcmdrepo "kali-auth-context/internal/infrastructure/db/repositories/tenant/commands"
	tenantqryrepo "kali-auth-context/internal/infrastructure/db/repositories/tenant/queries"
	usercmdrepo "kali-auth-context/internal/infrastructure/db/repositories/user/commands"
	userqryrepo "kali-auth-context/internal/infrastructure/db/repositories/user/queries"
	authhttpcommands "kali-auth-context/internal/infrastructure/http/contexts/auth/commands"
	authorizationhttpqueries "kali-auth-context/internal/infrastructure/http/contexts/authorization/queries"
	rbachttpcommands "kali-auth-context/internal/infrastructure/http/contexts/rbac/commands"
	rbachttpqueries "kali-auth-context/internal/infrastructure/http/contexts/rbac/queries"
	tenanthttpcommands "kali-auth-context/internal/infrastructure/http/contexts/tenant/commands"
	tenanthttpqueries "kali-auth-context/internal/infrastructure/http/contexts/tenant/queries"
	userhttpcommands "kali-auth-context/internal/infrastructure/http/contexts/user/commands"
	userhttpqueries "kali-auth-context/internal/infrastructure/http/contexts/user/queries"
	"kali-auth-context/internal/infrastructure/http/middleware"
	"kali-auth-context/internal/infrastructure/http/routes"
	"kali-auth-context/internal/infrastructure/providers"
	"kali-auth-context/internal/infrastructure/security"
)

var IdentitySet = wire.NewSet(
	InfrastructureSet,
	AccessApplicationSet,
	UserApplicationSet,
	TenantApplicationSet,
	RBACApplicationSet,
	NewSecurityModule,
	NewCoreContainer,
)

func InitializeContainer() (*CoreContainer, error) {
	cfg := config.LoadConfig()
	database, err := db.NewPool(cfg)
	if err != nil {
		return nil, err
	}

	migrator := db.NewMigrator(database)
	if err := migrator.Migrate(context.Background()); err != nil {
		return nil, err
	}

	uuidProvider := providers.NewUUIDProvider()
	passwordHasher := security.NewArgon2idHasher()
	accessTokenIssuer := security.NewAccessTokenIssuer(cfg)

	createUserRepo := usercmdrepo.NewCreateUserCommandRepository(database)
	updateUserRepo := usercmdrepo.NewUpdateUserCommandRepository(database)
	deleteUserRepo := usercmdrepo.NewDeleteUserCommandRepository(database)
	getUserByIdRepo := userqryrepo.NewGetUserByIdQueryRepository(database)
	getUserByEmailRepo := userqryrepo.NewGetUserByEmailQueryRepository(database)
	listUsersRepo := userqryrepo.NewListUsersQueryRepository(database)

	createTenantRepo := tenantcmdrepo.NewCreateTenantCommandRepository(database)
	updateTenantRepo := tenantcmdrepo.NewUpdateTenantCommandRepository(database)
	activateTenantRepo := tenantcmdrepo.NewActivateTenantCommandRepository(database)
	suspendTenantRepo := tenantcmdrepo.NewSuspendTenantCommandRepository(database)
	getTenantByIdRepo := tenantqryrepo.NewGetTenantByIdQueryRepository(database)
	getTenantByNameRepo := tenantqryrepo.NewGetTenantByNameQueryRepository(database)
	listTenantsRepo := tenantqryrepo.NewListTenantsQueryRepository(database)

	createRoleRepo := rbaccmdrepo.NewCreateRoleCommandRepository(database)
	updateRoleRepo := rbaccmdrepo.NewUpdateRoleCommandRepository(database)
	deleteRoleRepo := rbaccmdrepo.NewDeleteRoleCommandRepository(database)
	createPermissionRepo := rbaccmdrepo.NewCreatePermissionCommandRepository(database)
	updatePermissionRepo := rbaccmdrepo.NewUpdatePermissionCommandRepository(database)
	deletePermissionRepo := rbaccmdrepo.NewDeletePermissionCommandRepository(database)
	assignRoleToUserRepo := rbaccmdrepo.NewAssignRoleToUserCommandRepository(database)
	removeRoleFromUserRepo := rbaccmdrepo.NewRemoveRoleFromUserCommandRepository(database)
	assignPermissionToRoleRepo := rbaccmdrepo.NewAssignPermissionToRoleCommandRepository(database)
	removePermissionFromRoleRepo := rbaccmdrepo.NewRemovePermissionFromRoleCommandRepository(database)
	getRoleByIdRepo := rbacqryrepo.NewGetRoleByIdQueryRepository(database)
	listRolesRepo := rbacqryrepo.NewListRolesQueryRepository(database)
	getPermissionByIdRepo := rbacqryrepo.NewGetPermissionByIdQueryRepository(database)
	listPermissionsRepo := rbacqryrepo.NewListPermissionsQueryRepository(database)
	getUserRolesRepo := rbacqryrepo.NewGetUserRolesQueryRepository(database)
	getRolePermissionsRepo := rbacqryrepo.NewGetRolePermissionsQueryRepository(database)

	idempotencyRepo := idempotencyrepo.NewIdempotencyRepository(database)
	idempotencyMW := middleware.NewIdempotencyMiddleware(idempotencyRepo)

	jwtAuthMW := middleware.NewJWTAuthMiddleware(middleware.JWTAuthConfig{
		Secret:    cfg.JWTSecret,
		Issuer:    cfg.JWTIssuer,
		Audience:  cfg.JWTAudience,
		ClockSkew: 30 * time.Second,
	})

	securityModule := NewSecurityModule(uuidProvider, passwordHasher)

	authModule := NewAuthModule(
		authcommands.NewLoginCommand(getUserByEmailRepo, getUserRolesRepo, getRoleByIdRepo, getRolePermissionsRepo, getPermissionByIdRepo, passwordHasher),
	)

	authorizationModule := NewAuthorizationModule(
		authorizationqueries.NewAuthorizeQuery(getTenantByIdRepo, getUserByIdRepo, getUserRolesRepo, getRolePermissionsRepo, getPermissionByIdRepo),
	)

	userModule := NewUserModule(
		NewUserCommands(
			usercommands.NewCreateUserCommand(createUserRepo, uuidProvider, passwordHasher),
			usercommands.NewUpdateUserCommand(updateUserRepo, passwordHasher),
			usercommands.NewDeleteUserCommand(deleteUserRepo),
		),
		NewUserQueries(
			userqueries.NewGetUserByIdQuery(getUserByIdRepo),
			userqueries.NewGetUserByEmailQuery(getUserByEmailRepo),
			userqueries.NewListUsersQuery(listUsersRepo),
		),
	)

	tenantModule := NewTenantModule(
		NewTenantCommands(
			tenantcommands.NewCreateTenantCommand(createTenantRepo, uuidProvider),
			tenantcommands.NewUpdateTenantCommand(updateTenantRepo),
			tenantcommands.NewActivateTenantCommand(activateTenantRepo),
			tenantcommands.NewSuspendTenantCommand(suspendTenantRepo),
		),
		NewTenantQueries(
			tenantqueries.NewGetTenantByIdQuery(getTenantByIdRepo),
			tenantqueries.NewGetTenantByNameQuery(getTenantByNameRepo),
			tenantqueries.NewListTenantsQuery(listTenantsRepo),
		),
	)

	rbacModule := NewRBACModule(
		NewRBACCommands(
			rbaccommands.NewCreateRoleCommand(createRoleRepo, uuidProvider),
			rbaccommands.NewUpdateRoleCommand(updateRoleRepo),
			rbaccommands.NewDeleteRoleCommand(deleteRoleRepo),
			rbaccommands.NewCreatePermissionCommand(createPermissionRepo, uuidProvider),
			rbaccommands.NewUpdatePermissionCommand(updatePermissionRepo),
			rbaccommands.NewDeletePermissionCommand(deletePermissionRepo),
			rbaccommands.NewAssignRoleToUserCommand(assignRoleToUserRepo),
			rbaccommands.NewRemoveRoleFromUserCommand(removeRoleFromUserRepo),
			rbaccommands.NewAssignPermissionToRoleCommand(assignPermissionToRoleRepo),
			rbaccommands.NewRemovePermissionFromRoleCommand(removePermissionFromRoleRepo),
		),
		NewRBACQueries(
			rbacqueries.NewGetRoleByIdQuery(getRoleByIdRepo),
			rbacqueries.NewListRolesQuery(listRolesRepo),
			rbacqueries.NewGetPermissionByIdQuery(getPermissionByIdRepo),
			rbacqueries.NewListPermissionsQuery(listPermissionsRepo),
			rbacqueries.NewGetUserRolesQuery(getUserRolesRepo),
			rbacqueries.NewGetRolePermissionsQuery(getRolePermissionsRepo),
			rbacqueries.NewGetUserEffectivePermissionsQuery(getUserRolesRepo, getRolePermissionsRepo, getPermissionByIdRepo),
		),
	)

	loginHandler := authhttpcommands.NewLoginHandler(authModule.Login, accessTokenIssuer)
	authorizeHandler := authorizationhttpqueries.NewCheckHandler(authorizationModule.Authorize)

	createUserHandler := userhttpcommands.NewCreateHandler(userModule.Commands.Create)
	updateUserHandler := userhttpcommands.NewUpdateHandler(userModule.Commands.Update)
	deleteUserHandler := userhttpcommands.NewDeleteHandler(userModule.Commands.Delete)
	getUserByIdHandler := userhttpqueries.NewGetByIdHandler(userModule.Queries.GetById)
	getUserByEmailHandler := userhttpqueries.NewGetByEmailHandler(userModule.Queries.GetByEmail)
	listUsersHandler := userhttpqueries.NewListHandler(userModule.Queries.List)

	createTenantHandler := tenanthttpcommands.NewCreateHandler(tenantModule.Commands.Create)
	updateTenantHandler := tenanthttpcommands.NewUpdateHandler(tenantModule.Commands.Update)
	activateTenantHandler := tenanthttpcommands.NewActivateHandler(tenantModule.Commands.Activate)
	suspendTenantHandler := tenanthttpcommands.NewSuspendHandler(tenantModule.Commands.Suspend)
	getTenantByIdHandler := tenanthttpqueries.NewGetByIdHandler(tenantModule.Queries.GetById)
	getTenantByNameHandler := tenanthttpqueries.NewGetByNameHandler(tenantModule.Queries.GetByName)
	listTenantsHandler := tenanthttpqueries.NewListHandler(tenantModule.Queries.List)

	createRoleHandler := rbachttpcommands.NewCreateRoleHandler(rbacModule.Commands.CreateRole)
	updateRoleHandler := rbachttpcommands.NewUpdateRoleHandler(rbacModule.Commands.UpdateRole)
	deleteRoleHandler := rbachttpcommands.NewDeleteRoleHandler(rbacModule.Commands.DeleteRole)
	createPermissionHandler := rbachttpcommands.NewCreatePermissionHandler(rbacModule.Commands.CreatePermission)
	updatePermissionHandler := rbachttpcommands.NewUpdatePermissionHandler(rbacModule.Commands.UpdatePermission)
	deletePermissionHandler := rbachttpcommands.NewDeletePermissionHandler(rbacModule.Commands.DeletePermission)
	assignRoleToUserHandler := rbachttpcommands.NewAssignRoleToUserHandler(rbacModule.Commands.AssignRoleToUser)
	removeRoleFromUserHandler := rbachttpcommands.NewRemoveRoleFromUserHandler(rbacModule.Commands.RemoveRoleFromUser)
	assignPermissionToRoleHandler := rbachttpcommands.NewAssignPermissionToRoleHandler(rbacModule.Commands.AssignPermissionToRole)
	removePermissionFromRoleHandler := rbachttpcommands.NewRemovePermissionFromRoleHandler(rbacModule.Commands.RemovePermissionFromRole)

	getRoleByIdHandler := rbachttpqueries.NewGetRoleByIdHandler(rbacModule.Queries.GetRoleById)
	listRolesHandler := rbachttpqueries.NewListRolesHandler(rbacModule.Queries.ListRoles)
	getPermissionByIdHandler := rbachttpqueries.NewGetPermissionByIdHandler(rbacModule.Queries.GetPermissionById)
	listPermissionsHandler := rbachttpqueries.NewListPermissionsHandler(rbacModule.Queries.ListPermissions)
	getUserRolesHandler := rbachttpqueries.NewGetUserRolesHandler(rbacModule.Queries.GetUserRoles)
	getRolePermissionsHandler := rbachttpqueries.NewGetRolePermissionsHandler(rbacModule.Queries.GetRolePermissions)
	getUserEffectivePermissionsHandler := rbachttpqueries.NewGetUserEffectivePermissionsHandler(rbacModule.Queries.GetUserEffectivePermissions)

	httpModule := &HTTPModule{
		Router: routes.NewRouter(
			jwtAuthMW,
			idempotencyMW,
			loginHandler,
			authorizeHandler,
			createUserHandler,
			updateUserHandler,
			deleteUserHandler,
			getUserByIdHandler,
			getUserByEmailHandler,
			listUsersHandler,
			createTenantHandler,
			updateTenantHandler,
			activateTenantHandler,
			suspendTenantHandler,
			getTenantByIdHandler,
			getTenantByNameHandler,
			listTenantsHandler,
			createRoleHandler,
			updateRoleHandler,
			deleteRoleHandler,
			createPermissionHandler,
			updatePermissionHandler,
			deletePermissionHandler,
			assignRoleToUserHandler,
			removeRoleFromUserHandler,
			assignPermissionToRoleHandler,
			removePermissionFromRoleHandler,
			getRoleByIdHandler,
			listRolesHandler,
			getPermissionByIdHandler,
			listPermissionsHandler,
			getUserRolesHandler,
			getRolePermissionsHandler,
			getUserEffectivePermissionsHandler,
		),
	}

	return NewCoreContainer(cfg, database, securityModule, authModule, authorizationModule, httpModule, userModule, tenantModule, rbacModule), nil
}
