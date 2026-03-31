package mappers

import (
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
)

func ToDomainTenant(model *models.TenantModel) *identity.Tenant {
	return &identity.Tenant{
		Id:     identity.TenantId(model.Id),
		Name:   model.Name,
		Status: identity.TenantStatus(model.Status),
	}
}

func ToTenantModel(tenant *identity.Tenant) *models.TenantModel {
	return &models.TenantModel{
		Id:     string(tenant.Id),
		Name:   tenant.Name,
		Status: string(tenant.Status),
	}
}
