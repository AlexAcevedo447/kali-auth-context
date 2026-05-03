package routes

import (
	authcommands "kali-auth-context/internal/infrastructure/http/contexts/auth/commands"
	authorizationqueries "kali-auth-context/internal/infrastructure/http/contexts/authorization/queries"
	rbaccommands "kali-auth-context/internal/infrastructure/http/contexts/rbac/commands"
	rbacqueries "kali-auth-context/internal/infrastructure/http/contexts/rbac/queries"
	tenantcommands "kali-auth-context/internal/infrastructure/http/contexts/tenant/commands"
	tenantqueries "kali-auth-context/internal/infrastructure/http/contexts/tenant/queries"
	usercommands "kali-auth-context/internal/infrastructure/http/contexts/user/commands"
	userqueries "kali-auth-context/internal/infrastructure/http/contexts/user/queries"

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

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public endpoints (no JWT required)
	v1.Post("/auth/login", r.login.Handle)

	// Protected endpoints (JWT required)
	protected := v1.Group("", r.jwtAuth)
	protected.Post("/auth/authorize", r.authorize.Handle)

	users := protected.Group("/users")
	users.Post("/", r.idempotency, r.createUser.Handle)
	users.Get("/by-email", r.getUserByEmail.Handle)
	users.Get("/", r.listUsers.Handle)
	users.Put("/:userId", r.idempotency, r.updateUser.Handle)
	users.Delete("/:userId", r.idempotency, r.deleteUser.Handle)
	users.Get("/:userId", r.getUserById.Handle)

	tenants := protected.Group("/tenants")
	tenants.Post("/", r.idempotency, r.createTenant.Handle)
	tenants.Get("/by-name", r.getTenantByName.Handle)
	tenants.Get("/", r.listTenants.Handle)
	tenants.Put("/:tenantId", r.idempotency, r.updateTenant.Handle)
	tenants.Post("/:tenantId/activate", r.idempotency, r.activateTenant.Handle)
	tenants.Post("/:tenantId/suspend", r.idempotency, r.suspendTenant.Handle)
	tenants.Get("/:tenantId", r.getTenantById.Handle)

	roles := protected.Group("/roles")
	roles.Post("/", r.idempotency, r.createRole.Handle)
	roles.Get("/", r.listRoles.Handle)
	roles.Put("/:roleId", r.idempotency, r.updateRole.Handle)
	roles.Delete("/:roleId", r.idempotency, r.deleteRole.Handle)
	roles.Get("/:roleId", r.getRoleById.Handle)

	permissions := protected.Group("/permissions")
	permissions.Post("/", r.idempotency, r.createPermission.Handle)
	permissions.Get("/", r.listPermissions.Handle)
	permissions.Put("/:permissionId", r.idempotency, r.updatePermission.Handle)
	permissions.Delete("/:permissionId", r.idempotency, r.deletePermission.Handle)
	permissions.Get("/:permissionId", r.getPermissionById.Handle)

	rbac := protected.Group("/rbac")
	rbac.Post("/users/roles/assign", r.idempotency, r.assignRoleToUser.Handle)
	rbac.Post("/users/roles/remove", r.removeRoleFromUser.Handle)
	rbac.Post("/roles/permissions/assign", r.idempotency, r.assignPermissionToRole.Handle)
	rbac.Post("/roles/permissions/remove", r.removePermissionFromRole.Handle)
	rbac.Get("/users/roles", r.getUserRoles.Handle)
	rbac.Get("/roles/permissions", r.getRolePermissions.Handle)
	rbac.Get("/users/effective-permissions", r.getUserEffectivePermission.Handle)
}
