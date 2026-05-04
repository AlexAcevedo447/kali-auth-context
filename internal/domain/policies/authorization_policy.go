package policies

import (
	"strings"

	"kali-auth-context/internal/domain/identity"
)

func HasPermission(tenantId identity.TenantId, permissions []*identity.Permission, resource string, action string) bool {
	normalizedResource := normalizeAuthorizationValue(resource)
	normalizedAction := normalizeAuthorizationValue(action)
	if tenantId == "" || normalizedResource == "" || normalizedAction == "" {
		return false
	}

	for _, permission := range permissions {
		if permission == nil {
			continue
		}
		if permission.TenantId != tenantId {
			continue
		}
		permissionResource := normalizeAuthorizationValue(permission.Resource)
		permissionAction := normalizeAuthorizationValue(permission.Action)

		resourceAllowed := permissionResource == "*" || permissionResource == normalizedResource
		actionAllowed := permissionAction == "*" || permissionAction == normalizedAction

		if !resourceAllowed {
			continue
		}
		if !actionAllowed {
			continue
		}
		return true
	}

	return false
}

func normalizeAuthorizationValue(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
