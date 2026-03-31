package identity

import "errors"

type TenantId string

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "ACTIVE"
	TenantStatusSuspended TenantStatus = "SUSPENDED"
)

var (
	ErrTenantRequired     = errors.New("tenant_id is required")
	ErrTenantNameRequired = errors.New("tenant name is required")
)

// Tenant representa la entidad de negocio dueña del contexto de usuarios.
type Tenant struct {
	Id     TenantId
	Name   string
	Status TenantStatus
}

func NewTenant(id TenantId, name string) (*Tenant, error) {
	if id == "" {
		return nil, ErrTenantRequired
	}

	if name == "" {
		return nil, ErrTenantNameRequired
	}

	return &Tenant{
		Id:     id,
		Name:   name,
		Status: TenantStatusActive,
	}, nil
}

func (t *Tenant) Suspend() {
	t.Status = TenantStatusSuspended
}

func (t *Tenant) Activate() {
	t.Status = TenantStatusActive
}

func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}
