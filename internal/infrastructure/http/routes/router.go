package routes

import (
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"
	authcommands "kali-auth-context/internal/infrastructure/http/contexts/auth/commands"
	authorizationqueries "kali-auth-context/internal/infrastructure/http/contexts/authorization/queries"
	rbaccommands "kali-auth-context/internal/infrastructure/http/contexts/rbac/commands"
	rbacqueries "kali-auth-context/internal/infrastructure/http/contexts/rbac/queries"
	tenantcommands "kali-auth-context/internal/infrastructure/http/contexts/tenant/commands"
	tenantqueries "kali-auth-context/internal/infrastructure/http/contexts/tenant/queries"
	usercommands "kali-auth-context/internal/infrastructure/http/contexts/user/commands"
	userqueries "kali-auth-context/internal/infrastructure/http/contexts/user/queries"

	"kali-auth-context/internal/infrastructure/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// idempotencyHandler is a Fiber middleware. It is stored as a field so it can
// be injected via the constructor and applied only to non-naturally-idempotent
// POST endpoints (resource creation and relationship assignment).
type idempotencyHandler = fiber.Handler

// jwtAuthHandler is a Fiber middleware for JWT validation.
type jwtAuthHandler = fiber.Handler

type Router struct {
	jwtAuth                    jwtAuthHandler
	idempotency                idempotencyHandler
	login                      *authcommands.LoginHandler
	authorize                  *authorizationqueries.CheckHandler
	createUser                 *usercommands.CreateHandler
	updateUser                 *usercommands.UpdateHandler
	deleteUser                 *usercommands.DeleteHandler
	getUserById                *userqueries.GetByIdHandler
	getUserByEmail             *userqueries.GetByEmailHandler
	listUsers                  *userqueries.ListHandler
	createTenant               *tenantcommands.CreateHandler
	updateTenant               *tenantcommands.UpdateHandler
	activateTenant             *tenantcommands.ActivateHandler
	suspendTenant              *tenantcommands.SuspendHandler
	getTenantById              *tenantqueries.GetByIdHandler
	getTenantByName            *tenantqueries.GetByNameHandler
	listTenants                *tenantqueries.ListHandler
	createRole                 *rbaccommands.CreateRoleHandler
	updateRole                 *rbaccommands.UpdateRoleHandler
	deleteRole                 *rbaccommands.DeleteRoleHandler
	createPermission           *rbaccommands.CreatePermissionHandler
	updatePermission           *rbaccommands.UpdatePermissionHandler
	deletePermission           *rbaccommands.DeletePermissionHandler
	assignRoleToUser           *rbaccommands.AssignRoleToUserHandler
	removeRoleFromUser         *rbaccommands.RemoveRoleFromUserHandler
	assignPermissionToRole     *rbaccommands.AssignPermissionToRoleHandler
	removePermissionFromRole   *rbaccommands.RemovePermissionFromRoleHandler
	getRoleById                *rbacqueries.GetRoleByIdHandler
	listRoles                  *rbacqueries.ListRolesHandler
	getPermissionById          *rbacqueries.GetPermissionByIdHandler
	listPermissions            *rbacqueries.ListPermissionsHandler
	getUserRoles               *rbacqueries.GetUserRolesHandler
	getRolePermissions         *rbacqueries.GetRolePermissionsHandler
	getUserEffectivePermission *rbacqueries.GetUserEffectivePermissionsHandler
	jwksHandler                fiber.Handler
	guards                     *middleware.TenantGuards
	rbacGuards                 *middleware.RBACGuards
}

// WithTenantGuards activa el aislamiento multi-tenant en las rutas protegidas.
func (r *Router) WithTenantGuards(guards *middleware.TenantGuards) *Router {
	r.guards = guards
	return r
}

// WithRBACGuards activa la autorización por operación (MATRIZ_PERMISOS.md).
func (r *Router) WithRBACGuards(guards *middleware.RBACGuards) *Router {
	r.rbacGuards = guards
	return r
}

func NewRouter(
	jwtAuth jwtAuthHandler,
	idempotency idempotencyHandler,
	login *authcommands.LoginHandler,
	authorize *authorizationqueries.CheckHandler,
	createUser *usercommands.CreateHandler,
	updateUser *usercommands.UpdateHandler,
	deleteUser *usercommands.DeleteHandler,
	getUserById *userqueries.GetByIdHandler,
	getUserByEmail *userqueries.GetByEmailHandler,
	listUsers *userqueries.ListHandler,
	createTenant *tenantcommands.CreateHandler,
	updateTenant *tenantcommands.UpdateHandler,
	activateTenant *tenantcommands.ActivateHandler,
	suspendTenant *tenantcommands.SuspendHandler,
	getTenantById *tenantqueries.GetByIdHandler,
	getTenantByName *tenantqueries.GetByNameHandler,
	listTenants *tenantqueries.ListHandler,
	createRole *rbaccommands.CreateRoleHandler,
	updateRole *rbaccommands.UpdateRoleHandler,
	deleteRole *rbaccommands.DeleteRoleHandler,
	createPermission *rbaccommands.CreatePermissionHandler,
	updatePermission *rbaccommands.UpdatePermissionHandler,
	deletePermission *rbaccommands.DeletePermissionHandler,
	assignRoleToUser *rbaccommands.AssignRoleToUserHandler,
	removeRoleFromUser *rbaccommands.RemoveRoleFromUserHandler,
	assignPermissionToRole *rbaccommands.AssignPermissionToRoleHandler,
	removePermissionFromRole *rbaccommands.RemovePermissionFromRoleHandler,
	getRoleById *rbacqueries.GetRoleByIdHandler,
	listRoles *rbacqueries.ListRolesHandler,
	getPermissionById *rbacqueries.GetPermissionByIdHandler,
	listPermissions *rbacqueries.ListPermissionsHandler,
	getUserRoles *rbacqueries.GetUserRolesHandler,
	getRolePermissions *rbacqueries.GetRolePermissionsHandler,
	getUserEffectivePermission *rbacqueries.GetUserEffectivePermissionsHandler,
	jwksHandler fiber.Handler,
) *Router {
	return &Router{
		jwtAuth:                    jwtAuth,
		idempotency:                idempotency,
		login:                      login,
		authorize:                  authorize,
		createUser:                 createUser,
		updateUser:                 updateUser,
		deleteUser:                 deleteUser,
		getUserById:                getUserById,
		getUserByEmail:             getUserByEmail,
		listUsers:                  listUsers,
		createTenant:               createTenant,
		updateTenant:               updateTenant,
		activateTenant:             activateTenant,
		suspendTenant:              suspendTenant,
		getTenantById:              getTenantById,
		getTenantByName:            getTenantByName,
		listTenants:                listTenants,
		createRole:                 createRole,
		updateRole:                 updateRole,
		deleteRole:                 deleteRole,
		createPermission:           createPermission,
		updatePermission:           updatePermission,
		deletePermission:           deletePermission,
		assignRoleToUser:           assignRoleToUser,
		removeRoleFromUser:         removeRoleFromUser,
		assignPermissionToRole:     assignPermissionToRole,
		removePermissionFromRole:   removePermissionFromRole,
		getRoleById:                getRoleById,
		listRoles:                  listRoles,
		getPermissionById:          getPermissionById,
		listPermissions:            listPermissions,
		getUserRoles:               getUserRoles,
		getRolePermissions:         getRolePermissions,
		getUserEffectivePermission: getUserEffectivePermission,
		jwksHandler:                jwksHandler,
	}
}

func (r *Router) Register(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "kali-auth-context",
			"status":  "up",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Public key set (JWKS) for services that verify RS256 access tokens.
	app.Get("/.well-known/jwks.json", r.jwksHandler)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public endpoints (no JWT required)
	v1.Post("/auth/login", r.login.Handle)

	// Protected endpoints (JWT required)
	// Todas las rutas protegidas operan en el tenant del token (ver TenantGuards).
	// ActiveTenant consulta el estado de la empresa en cada petición: suspenderla
	// bloquea también los tokens ya emitidos.
	protected := v1.Group("", r.jwtAuth, r.guards.ActiveTenant(), r.guards.Scope())
	ownedUser := middleware.OwnedRef{Source: "param", Field: "userId", Table: "user_models"}
	ownedRole := middleware.OwnedRef{Source: "param", Field: "roleId", Table: "roles"}
	ownedPermission := middleware.OwnedRef{Source: "param", Field: "permissionId", Table: "permissions"}
	bodyUser := middleware.OwnedRef{Source: "body", Field: "user_id", Table: "user_models"}
	bodyRole := middleware.OwnedRef{Source: "body", Field: "role_id", Table: "roles"}
	bodyPermission := middleware.OwnedRef{Source: "body", Field: "permission_id", Table: "permissions"}
	queryUser := middleware.OwnedRef{Source: "query", Field: "user_id", Table: "user_models"}
	queryRole := middleware.OwnedRef{Source: "query", Field: "role_id", Table: "roles"}
	rb := r.rbacGuards
	// Evaluar los permisos propios no exige nada más; los de otro usuario, rbac/manage.
	protected.Post("/auth/authorize", rb.SelfOr("body", "user_id", "rbac", "manage"), r.authorize.Handle)
	// Comprobación de sesión para otros servicios (p. ej. facturación): responde 200
	// solo si el token es válido y su empresa está activa (403 tenant is suspended),
	// con los permisos efectivos del propio usuario.
	protected.Get("/auth/session", rb.SessionHandler)

	users := protected.Group("/users")
	// users/read para consultar (manage incluye read) y users/manage para modificar.
	// Editar o eliminar a otro usuario exige cubrir sus permisos; eliminar no puede
	// dejar a la empresa sin administrador efectivo.
	users.Post("/", rb.Require("users", "manage"), r.idempotency, r.createUser.Handle)
	users.Get("/by-email", rb.Require("users", "read"), r.getUserByEmail.Handle)
	users.Get("/", rb.RequireAny(grant("users", "read"), grant("rbac", "manage")), r.listUsers.Handle)
	users.Put("/:userId", r.guards.Owned(ownedUser), rb.Require("users", "manage"), rb.ManageUserTarget("param", "userId"), r.idempotency, r.updateUser.Handle)
	users.Delete("/:userId", r.guards.Owned(ownedUser), rb.Require("users", "manage"), rb.NotSelf("param", "userId", identity.ErrSelfDeletion), rb.ManageUserTarget("param", "userId"), rb.KeepAdministrator(middleware.ChangeDeleteUser), r.idempotency, r.deleteUser.Handle)
	users.Get("/:userId", r.guards.Owned(ownedUser), rb.SelfOr("param", "userId", "users", "read"), r.getUserById.Handle)

	// Administración de empresas: solo el tenant de plataforma.
	tenants := protected.Group("/tenants", r.guards.PlatformAdmin())
	tenants.Post("/", r.idempotency, r.createTenant.Handle)
	tenants.Get("/by-name", r.getTenantByName.Handle)
	tenants.Get("/", r.listTenants.Handle)
	tenants.Put("/:tenantId", r.idempotency, r.updateTenant.Handle)
	tenants.Post("/:tenantId/activate", r.idempotency, r.activateTenant.Handle)
	tenants.Post("/:tenantId/suspend", r.idempotency, r.suspendTenant.Handle)
	tenants.Get("/:tenantId", r.getTenantById.Handle)

	roles := protected.Group("/roles")
	roles.Post("/", rb.Require("roles", "manage"), r.idempotency, r.createRole.Handle)
	roles.Get("/", rb.RequireAny(grant("roles", "read"), grant("rbac", "manage")), r.listRoles.Handle)
	roles.Put("/:roleId", r.guards.Owned(ownedRole), rb.Require("roles", "manage"), rb.ProtectSeedRole("param", "roleId"), r.idempotency, r.updateRole.Handle)
	roles.Delete("/:roleId", r.guards.Owned(ownedRole), rb.Require("roles", "manage"), rb.ProtectSeedRole("param", "roleId"), r.idempotency, r.deleteRole.Handle)
	roles.Get("/:roleId", r.guards.Owned(ownedRole), rb.RequireAny(grant("roles", "read"), grant("rbac", "manage")), r.getRoleById.Handle)

	permissions := protected.Group("/permissions")
	permissions.Get("/catalog", middleware.CatalogHandler)
	permissions.Post("/", rb.Require("permissions", "manage"), rb.NoEscalationPermissionWrite(), r.idempotency, r.createPermission.Handle)
	permissions.Get("/", rb.RequireAny(grant("permissions", "read"), grant("rbac", "manage")), r.listPermissions.Handle)
	permissions.Put("/:permissionId", r.guards.Owned(ownedPermission), rb.Require("permissions", "manage"), rb.ProtectSeedPermission("param", "permissionId"), rb.NoEscalationPermissionWrite(), r.idempotency, r.updatePermission.Handle)
	permissions.Delete("/:permissionId", r.guards.Owned(ownedPermission), rb.Require("permissions", "manage"), rb.ProtectSeedPermission("param", "permissionId"), r.idempotency, r.deletePermission.Handle)
	permissions.Get("/:permissionId", r.guards.Owned(ownedPermission), rb.RequireAny(grant("permissions", "read"), grant("rbac", "manage")), r.getPermissionById.Handle)

	rbac := protected.Group("/rbac")
	// Usuario, rol y permiso deben ser del mismo tenant (Owned); además rbac/manage,
	// no escalar privilegios y no dejar la empresa sin administrador efectivo.
	rbac.Post("/users/roles/assign", r.guards.Owned(bodyUser, bodyRole), rb.Require("rbac", "manage"), rb.NoEscalationUserRole(), r.idempotency, r.assignRoleToUser.Handle)
	rbac.Post("/users/roles/remove", r.guards.Owned(bodyUser, bodyRole), rb.Require("rbac", "manage"), rb.NoEscalationUserRole(), rb.KeepAdministrator(middleware.ChangeRemoveUserRole), r.removeRoleFromUser.Handle)
	rbac.Post("/roles/permissions/assign", r.guards.Owned(bodyRole, bodyPermission), rb.Require("rbac", "manage"), rb.NoEscalationRolePermission(), r.idempotency, r.assignPermissionToRole.Handle)
	rbac.Post("/roles/permissions/remove", r.guards.Owned(bodyRole, bodyPermission), rb.Require("rbac", "manage"), rb.ProtectSeedRole("body", "role_id"), rb.NoEscalationRolePermission(), rb.KeepAdministrator(middleware.ChangeRemoveRolePermission), r.removePermissionFromRole.Handle)
	rbac.Get("/users/roles", r.guards.Owned(queryUser), rb.SelfOr("query", "user_id", "rbac", "manage"), r.getUserRoles.Handle)
	rbac.Get("/roles/permissions", r.guards.Owned(queryRole), rb.RequireAny(grant("roles", "read"), grant("rbac", "manage")), r.getRolePermissions.Handle)
	rbac.Get("/users/effective-permissions", r.guards.Owned(queryUser), rb.SelfOr("query", "user_id", "rbac", "manage"), r.getUserEffectivePermission.Handle)
}

// grant abrevia un permiso (recurso, acción) para RequireAny.
func grant(resource, action string) policies.Grant {
	return policies.Grant{Resource: resource, Action: action}
}
