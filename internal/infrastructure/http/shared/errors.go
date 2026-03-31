package shared

import (
	"database/sql"
	"errors"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"

	"github.com/gofiber/fiber/v2"
)

func WriteError(c *fiber.Ctx, err error) error {
	if err == nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "resource not found"})
	}

	if errors.Is(err, identity.ErrInvalidCredentials) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if isValidationError(err) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}

func isValidationError(err error) bool {
	return errors.Is(err, identity.ErrTenantRequired) ||
		errors.Is(err, identity.ErrTenantNameRequired) ||
		errors.Is(err, identity.ErrUserIdRequired) ||
		errors.Is(err, identity.ErrRoleIdRequired) ||
		errors.Is(err, identity.ErrPermissionIdRequired) ||
		errors.Is(err, identity.ErrEmailRequired) ||
		errors.Is(err, identity.ErrPasswordRequired) ||
		errors.Is(err, identity.ErrRoleNameRequired) ||
		errors.Is(err, identity.ErrPermissionResource) ||
		errors.Is(err, identity.ErrPermissionAction) ||
		errors.Is(err, identity.ErrAuthorizationRequestRequired) ||
		errors.Is(err, identity.ErrAuthorizationResourceRequired) ||
		errors.Is(err, identity.ErrAuthorizationActionRequired) ||
		errors.Is(err, policies.ErrPasswordTooShort) ||
		errors.Is(err, policies.ErrPasswordMissingUpper) ||
		errors.Is(err, policies.ErrPasswordMissingLower) ||
		errors.Is(err, policies.ErrPasswordMissingDigit) ||
		errors.Is(err, policies.ErrPasswordMissingSymbol) ||
		errors.Is(err, policies.ErrPasswordHasWhitespace)
}
