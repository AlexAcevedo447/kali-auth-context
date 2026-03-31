package mappers

import (
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
)

func ToRoleModel(role *identity.Role) *models.RoleModel {
	return &models.RoleModel{
		Id:          string(role.Id),
		TenantId:    string(role.TenantId),
		Name:        role.Name,
		Description: role.Description,
	}
}

func ToDomainRole(model *models.RoleModel) *identity.Role {
	return &identity.Role{
		Id:          identity.RoleId(model.Id),
		TenantId:    identity.TenantId(model.TenantId),
		Name:        model.Name,
		Description: model.Description,
	}
}

func ToPermissionModel(permission *identity.Permission) *models.PermissionModel {
	return &models.PermissionModel{
		Id:       string(permission.Id),
		TenantId: string(permission.TenantId),
		Resource: permission.Resource,
		Action:   permission.Action,
	}
}

func ToDomainPermission(model *models.PermissionModel) *identity.Permission {
	return &identity.Permission{
		Id:       identity.PermissionId(model.Id),
		TenantId: identity.TenantId(model.TenantId),
		Resource: model.Resource,
		Action:   model.Action,
	}
}

func ToUserRoleModel(relation *identity.UserRole) *models.UserRoleModel {
	return &models.UserRoleModel{
		TenantId: string(relation.TenantId),
		UserId:   string(relation.UserId),
		RoleId:   string(relation.RoleId),
	}
}

func ToDomainUserRole(model *models.UserRoleModel) *identity.UserRole {
	return &identity.UserRole{
		TenantId: identity.TenantId(model.TenantId),
		UserId:   identity.UserId(model.UserId),
		RoleId:   identity.RoleId(model.RoleId),
	}
}

func ToRolePermissionModel(relation *identity.RolePermission) *models.RolePermissionModel {
	return &models.RolePermissionModel{
		TenantId:     string(relation.TenantId),
		RoleId:       string(relation.RoleId),
		PermissionId: string(relation.PermissionId),
	}
}

func ToDomainRolePermission(model *models.RolePermissionModel) *identity.RolePermission {
	return &identity.RolePermission{
		TenantId:     identity.TenantId(model.TenantId),
		RoleId:       identity.RoleId(model.RoleId),
		PermissionId: identity.PermissionId(model.PermissionId),
	}
}
