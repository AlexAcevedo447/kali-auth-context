package bootstrap

import (
	"kali-auth-context/internal/infrastructure/config"
	"kali-auth-context/internal/infrastructure/db"
	rbaccmdrepo "kali-auth-context/internal/infrastructure/db/repositories/rbac/commands"
	rbacqryrepo "kali-auth-context/internal/infrastructure/db/repositories/rbac/queries"
	tenantcmdrepo "kali-auth-context/internal/infrastructure/db/repositories/tenant/commands"
	tenantqryrepo "kali-auth-context/internal/infrastructure/db/repositories/tenant/queries"
	usercmdrepo "kali-auth-context/internal/infrastructure/db/repositories/user/commands"
	userqryrepo "kali-auth-context/internal/infrastructure/db/repositories/user/queries"
	"kali-auth-context/internal/infrastructure/providers"
	"kali-auth-context/internal/infrastructure/security"
	"kali-auth-context/internal/ports"

	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(config.LoadConfig)

var DatabaseSet = wire.NewSet(db.NewPool)

var ProviderSet = wire.NewSet(
	providers.NewUUIDProvider,
	wire.Bind(new(ports.IUUIDProvider), new(*providers.UUIDProvider)),
)

var SecuritySet = wire.NewSet(
	security.NewArgon2idHasher,
	wire.Bind(new(ports.IPasswordHasher), new(*security.Argon2idHasher)),
)

var UserRepositorySet = wire.NewSet(
	usercmdrepo.NewCreateUserCommandRepository,
	wire.Bind(new(ports.ICreateUserCommandRepository), new(*usercmdrepo.CreateUserCommandRepository)),
	usercmdrepo.NewUpdateUserCommandRepository,
	wire.Bind(new(ports.IUpdateUserCommandRepository), new(*usercmdrepo.UpdateUserCommandRepository)),
	usercmdrepo.NewDeleteUserCommandRepository,
	wire.Bind(new(ports.IDeleteUserCommandRepository), new(*usercmdrepo.DeleteUserCommandRepository)),
	userqryrepo.NewGetUserByIdQueryRepository,
	wire.Bind(new(ports.IGetUserByIdQueryRepository), new(*userqryrepo.GetUserByIdQueryRepository)),
	userqryrepo.NewGetUserByEmailQueryRepository,
	wire.Bind(new(ports.IGetUserByEmailQueryRepository), new(*userqryrepo.GetUserByEmailQueryRepository)),
	userqryrepo.NewListUsersQueryRepository,
	wire.Bind(new(ports.IListUsersQueryRepository), new(*userqryrepo.ListUsersQueryRepository)),
)

var TenantRepositorySet = wire.NewSet(
	tenantcmdrepo.NewCreateTenantCommandRepository,
	wire.Bind(new(ports.ICreateTenantCommandRepository), new(*tenantcmdrepo.CreateTenantCommandRepository)),
	tenantcmdrepo.NewUpdateTenantCommandRepository,
	wire.Bind(new(ports.IUpdateTenantCommandRepository), new(*tenantcmdrepo.UpdateTenantCommandRepository)),
	tenantcmdrepo.NewActivateTenantCommandRepository,
	wire.Bind(new(ports.IActivateTenantCommandRepository), new(*tenantcmdrepo.ActivateTenantCommandRepository)),
	tenantcmdrepo.NewSuspendTenantCommandRepository,
	wire.Bind(new(ports.ISuspendTenantCommandRepository), new(*tenantcmdrepo.SuspendTenantCommandRepository)),
	tenantqryrepo.NewGetTenantByIdQueryRepository,
	wire.Bind(new(ports.IGetTenantByIdQueryRepository), new(*tenantqryrepo.GetTenantByIdQueryRepository)),
	tenantqryrepo.NewGetTenantByNameQueryRepository,
	wire.Bind(new(ports.IGetTenantByNameQueryRepository), new(*tenantqryrepo.GetTenantByNameQueryRepository)),
	tenantqryrepo.NewListTenantsQueryRepository,
	wire.Bind(new(ports.IListTenantsQueryRepository), new(*tenantqryrepo.ListTenantsQueryRepository)),
)

var RBACRepositorySet = wire.NewSet(
	rbaccmdrepo.NewCreateRoleCommandRepository,
	wire.Bind(new(ports.ICreateRoleCommandRepository), new(*rbaccmdrepo.CreateRoleCommandRepository)),
	rbaccmdrepo.NewUpdateRoleCommandRepository,
	wire.Bind(new(ports.IUpdateRoleCommandRepository), new(*rbaccmdrepo.UpdateRoleCommandRepository)),
	rbaccmdrepo.NewDeleteRoleCommandRepository,
	wire.Bind(new(ports.IDeleteRoleCommandRepository), new(*rbaccmdrepo.DeleteRoleCommandRepository)),
	rbaccmdrepo.NewCreatePermissionCommandRepository,
	wire.Bind(new(ports.ICreatePermissionCommandRepository), new(*rbaccmdrepo.CreatePermissionCommandRepository)),
	rbaccmdrepo.NewUpdatePermissionCommandRepository,
	wire.Bind(new(ports.IUpdatePermissionCommandRepository), new(*rbaccmdrepo.UpdatePermissionCommandRepository)),
	rbaccmdrepo.NewDeletePermissionCommandRepository,
	wire.Bind(new(ports.IDeletePermissionCommandRepository), new(*rbaccmdrepo.DeletePermissionCommandRepository)),
	rbaccmdrepo.NewAssignRoleToUserCommandRepository,
	wire.Bind(new(ports.IAssignRoleToUserCommandRepository), new(*rbaccmdrepo.AssignRoleToUserCommandRepository)),
	rbaccmdrepo.NewRemoveRoleFromUserCommandRepository,
	wire.Bind(new(ports.IRemoveRoleFromUserCommandRepository), new(*rbaccmdrepo.RemoveRoleFromUserCommandRepository)),
	rbaccmdrepo.NewAssignPermissionToRoleCommandRepository,
	wire.Bind(new(ports.IAssignPermissionToRoleCommandRepository), new(*rbaccmdrepo.AssignPermissionToRoleCommandRepository)),
	rbaccmdrepo.NewRemovePermissionFromRoleCommandRepository,
	wire.Bind(new(ports.IRemovePermissionFromRoleCommandRepository), new(*rbaccmdrepo.RemovePermissionFromRoleCommandRepository)),
	rbacqryrepo.NewGetRoleByIdQueryRepository,
	wire.Bind(new(ports.IGetRoleByIdQueryRepository), new(*rbacqryrepo.GetRoleByIdQueryRepository)),
	rbacqryrepo.NewListRolesQueryRepository,
	wire.Bind(new(ports.IListRolesQueryRepository), new(*rbacqryrepo.ListRolesQueryRepository)),
	rbacqryrepo.NewGetPermissionByIdQueryRepository,
	wire.Bind(new(ports.IGetPermissionByIdQueryRepository), new(*rbacqryrepo.GetPermissionByIdQueryRepository)),
	rbacqryrepo.NewListPermissionsQueryRepository,
	wire.Bind(new(ports.IListPermissionsQueryRepository), new(*rbacqryrepo.ListPermissionsQueryRepository)),
	rbacqryrepo.NewGetUserRolesQueryRepository,
	wire.Bind(new(ports.IGetUserRolesQueryRepository), new(*rbacqryrepo.GetUserRolesQueryRepository)),
	rbacqryrepo.NewGetRolePermissionsQueryRepository,
	wire.Bind(new(ports.IGetRolePermissionsQueryRepository), new(*rbacqryrepo.GetRolePermissionsQueryRepository)),
)

var InfrastructureSet = wire.NewSet(
	ConfigSet,
	DatabaseSet,
	ProviderSet,
	SecuritySet,
	UserRepositorySet,
	TenantRepositorySet,
	RBACRepositorySet,
)
