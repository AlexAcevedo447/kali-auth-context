package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

// LoginCandidate es una cuenta que coincide con el correo del login, junto con
// los datos de su empresa (tenant) necesarios para resolverla.
type LoginCandidate struct {
	User         *identity.User
	TenantName   string
	TenantCode   string
	TenantStatus identity.TenantStatus
}

// IFindLoginCandidatesQueryRepository busca las cuentas de un correo en todas las
// empresas o, si se indica tenantCode, solo en la empresa con ese código legible.
type IFindLoginCandidatesQueryRepository interface {
	FindByEmail(ctx context.Context, email, tenantCode string) ([]LoginCandidate, error)
}
