package identity

import (
	"errors"
	"regexp"
	"strings"
)

type TenantId string

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "ACTIVE"
	TenantStatusSuspended TenantStatus = "SUSPENDED"
)

var (
	ErrTenantRequired     = errors.New("tenant_id is required")
	ErrTenantNameRequired = errors.New("tenant name is required")
	ErrTenantCodeInvalid  = errors.New("tenant code must have 3-40 lowercase letters, digits or hyphens")
	ErrTenantCodeTaken    = errors.New("tenant code already exists")
	// ErrTenantCodeRequired indica que el correo y la contraseña corresponden a
	// cuentas de varias empresas: el usuario debe indicar el código de empresa.
	ErrTenantCodeRequired = errors.New("company code required")
	ErrTenantSuspended    = errors.New("tenant is suspended")
	// ErrTenantStatusInvalid: el estado solo puede ser ACTIVE o SUSPENDED.
	ErrTenantStatusInvalid = errors.New("tenant status must be ACTIVE or SUSPENDED")
	// ErrPlatformTenantSuspension: la empresa de plataforma no se puede suspender
	// (dejaría sin acceso a la administración de empresas).
	ErrPlatformTenantSuspension = errors.New("platform tenant cannot be suspended")
)

// tenantCodePattern define el código legible de la empresa (slug): es lo que
// una persona puede escribir para identificar su empresa; nunca el UUID.
var tenantCodePattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{1,38}[a-z0-9])$`)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

// SlugifyTenantCode deriva un código legible a partir del nombre de la empresa.
func SlugifyTenantCode(name string) string {
	replacer := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n")
	slug := nonSlugChars.ReplaceAllString(replacer.Replace(strings.ToLower(name)), "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 40 {
		slug = strings.Trim(slug[:40], "-")
	}
	for len(slug) < 3 {
		slug += "0"
	}
	return slug
}

// ValidateTenantCode comprueba el formato del código legible de la empresa.
func ValidateTenantCode(code string) error {
	if !tenantCodePattern.MatchString(code) {
		return ErrTenantCodeInvalid
	}
	return nil
}

// Tenant representa la entidad de negocio dueña del contexto de usuarios.
type Tenant struct {
	Id   TenantId
	Name string
	// Code es el identificador legible de la empresa (único). Solo se pide en el
	// login cuando un mismo correo y contraseña existen en varias empresas.
	Code   string
	Status TenantStatus
}

func NewTenant(id TenantId, name string) (*Tenant, error) {
	if id == "" {
		return nil, ErrTenantRequired
	}

	// Un nombre solo con espacios equivale a vacío; se guarda sin espacios a los lados.
	name = strings.TrimSpace(name)
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
