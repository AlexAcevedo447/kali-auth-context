package policies

import (
	"sort"
	"strings"

	"kali-auth-context/internal/domain/identity"
)

// Acciones y comodín del catálogo RBAC (ver docs/requisitos/MATRIZ_PERMISOS.md
// en invoice-kali-front).
const (
	Wildcard     = "*"
	ActionRead   = "read"
	ActionManage = "manage"
)

// Catalog: recursos y acciones válidos, basados en las operaciones existentes.
var Catalog = map[string][]string{
	"tenants":       {ActionManage},
	"users":         {ActionRead, ActionManage},
	"roles":         {ActionRead, ActionManage},
	"permissions":   {ActionRead, ActionManage},
	"rbac":          {ActionManage},
	"invoices":      {ActionRead, ActionManage},
	"invoice-items": {ActionRead, ActionManage},
	"customers":     {ActionRead, ActionManage},
	"products":      {ActionRead, ActionManage},
}

// CatalogEntry es la forma pública del catálogo.
type CatalogEntry struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

// CatalogEntries devuelve el catálogo ordenado por recurso.
func CatalogEntries() []CatalogEntry {
	out := make([]CatalogEntry, 0, len(Catalog))
	for resource, actions := range Catalog {
		out = append(out, CatalogEntry{Resource: resource, Actions: append([]string(nil), actions...)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Resource < out[j].Resource })
	return out
}

// NormalizePermissionValue: sin espacios a los lados y en minúsculas.
func NormalizePermissionValue(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

// ValidatePermission normaliza y valida (recurso, acción) contra el catálogo.
// "*" vale como recurso o acción; con recurso "*" la acción debe ser "*",
// "read" o "manage".
func ValidatePermission(resource, action string) (string, string, error) {
	r := NormalizePermissionValue(resource)
	a := NormalizePermissionValue(action)
	if r == "" {
		return "", "", identity.ErrPermissionResource
	}
	if a == "" {
		return "", "", identity.ErrPermissionAction
	}
	if r == Wildcard {
		if a != Wildcard && a != ActionRead && a != ActionManage {
			return "", "", identity.ErrUnknownAction
		}
		return r, a, nil
	}
	actions, ok := Catalog[r]
	if !ok {
		return "", "", identity.ErrUnknownResource
	}
	if a == Wildcard {
		return r, a, nil
	}
	for _, valid := range actions {
		if a == valid {
			return r, a, nil
		}
	}
	return "", "", identity.ErrUnknownAction
}

// Grant es un permiso (recurso, acción) sin identidad, para comparar alcances.
type Grant struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

func (g Grant) normalized() Grant {
	return Grant{Resource: NormalizePermissionValue(g.Resource), Action: NormalizePermissionValue(g.Action)}
}

// grantCovers indica si el permiso que se tiene (have) alcanza al pedido (want)
// considerando comodines y que manage incluye read.
func grantCovers(have, want Grant) bool {
	have, want = have.normalized(), want.normalized()
	resourceOK := have.Resource == Wildcard || have.Resource == want.Resource
	actionOK := have.Action == Wildcard || have.Action == want.Action ||
		(want.Action == ActionRead && have.Action == ActionManage)
	return resourceOK && actionOK
}

// Covers: alguno de los permisos efectivos alcanza a want. Un comodín solo se
// alcanza con un comodín igual o mayor (p. ej. */* solo lo cubre */*).
func Covers(effective []Grant, want Grant) bool {
	for _, have := range effective {
		if grantCovers(have, want) {
			return true
		}
	}
	return false
}

// CoversAll: alcanza todos los permisos pedidos (regla de no escalamiento).
func CoversAll(effective []Grant, wants []Grant) bool {
	for _, w := range wants {
		if !Covers(effective, w) {
			return false
		}
	}
	return true
}

// AdministratorGrants: lo que hace a un usuario "administrador efectivo" de su
// empresa (gestionar roles, permisos y asignaciones).
var AdministratorGrants = []Grant{
	{Resource: "roles", Action: ActionManage},
	{Resource: "permissions", Action: ActionManage},
	{Resource: "rbac", Action: ActionManage},
}

func IsEffectiveAdministrator(effective []Grant) bool {
	return CoversAll(effective, AdministratorGrants)
}

// IsSensitiveGrant: permisos que solo puede crear u otorgar quien ya los cubre
// (comodines y administración de plataforma).
func IsSensitiveGrant(g Grant) bool {
	g = g.normalized()
	return g.Resource == Wildcard || g.Action == Wildcard || g.Resource == "tenants"
}
