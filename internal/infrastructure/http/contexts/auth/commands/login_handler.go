package commands

import (
	"strings"

	authcommands "kali-auth-context/internal/application/auth/commands"
	"kali-auth-context/internal/infrastructure/http/shared"
	"kali-auth-context/internal/infrastructure/security"

	"github.com/gofiber/fiber/v2"
)

// Handler encargado de manejar el inicio de sesión (login) de usuarios.
// Recibe las credenciales desde la petición HTTP y ejecuta el caso de uso de autenticación.
type LoginHandler struct {
	command     *authcommands.LoginCommand
	tokenIssuer *security.AccessTokenIssuer
	// platformTenantID es la empresa de plataforma (SEED_MASTER_TENANT_ID), la
	// misma que usa middleware.TenantGuards.PlatformOnly para /tenants.
	platformTenantID string
}

func NewLoginHandler(command *authcommands.LoginCommand, tokenIssuer *security.AccessTokenIssuer, platformTenantID string) *LoginHandler {
	return &LoginHandler{command: command, tokenIssuer: tokenIssuer, platformTenantID: platformTenantID}
}

// loginRequest no admite tenant_id: el tenant se resuelve en el servidor. El
// código de empresa (legible) solo se usa si el correo existe en varias empresas.
type loginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	TenantCode string `json:"tenant_code"`
}

type loginResponse struct {
	AccessToken string              `json:"access_token"`
	TokenType   string              `json:"token_type"`
	ExpiresIn   int64               `json:"expires_in"`
	TenantId    string              `json:"tenant_id"`
	TenantName  string              `json:"tenant_name"`
	TenantCode  string              `json:"tenant_code"`
	UserId      string              `json:"user_id"`
	Email       string              `json:"email"`
	NeedsRehash bool                `json:"needs_rehash"`
	Roles       []string            `json:"roles"`
	Permissions []map[string]string `json:"permissions"`
	// IsPlatformTenant indica si la sesión pertenece a la empresa de plataforma,
	// la única autorizada a administrar empresas (/api/v1/tenants). Es solo
	// informativo para la interfaz: la autorización la aplica el servidor.
	IsPlatformTenant bool `json:"is_platform_tenant"`
}

// Maneja la petición POST para iniciar sesión.
// Recibe tenant_id, email y password en el cuerpo JSON.
// Si la autenticación es correcta, retorna un token de acceso y datos del usuario.
// Si las credenciales son incorrectas, retorna un error de autenticación.
func (h *LoginHandler) Handle(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Ejecuta el caso de uso de autenticación con las credenciales recibidas.
	result, err := h.command.Execute(&authcommands.LoginDto{
		Email:      strings.TrimSpace(req.Email),
		Password:   req.Password,
		TenantCode: strings.ToLower(strings.TrimSpace(req.TenantCode)),
	})
	if err != nil {
		// Si las credenciales son incorrectas, retorna error de autenticación.
		return shared.WriteError(c, err)
	}

	// Si la autenticación es correcta, genera el token de acceso.
	token, expiresIn, err := h.tokenIssuer.Issue(result.TenantId, result.UserId, result.Email, result.Roles, result.Permissions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate access token"})
	}

	// Prepara los datos de roles y permisos para la respuesta.
	roleNames := make([]string, 0, len(result.Roles))
	for _, r := range result.Roles {
		roleNames = append(roleNames, r.Name)
	}

	perms := make([]map[string]string, 0, len(result.Permissions))
	for _, p := range result.Permissions {
		perms = append(perms, map[string]string{
			"id":       string(p.Id),
			"resource": p.Resource,
			"action":   p.Action,
		})
	}

	// Devuelve respuesta de autenticación satisfactoria con el token y datos del usuario.
	return c.Status(fiber.StatusOK).JSON(loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		TenantId:    string(result.TenantId),
		TenantName:  result.TenantName,
		TenantCode:  result.TenantCode,
		UserId:      string(result.UserId),
		Email:       result.Email,
		NeedsRehash: result.NeedsRehash,
		Roles:       roleNames,
		Permissions: perms,
		IsPlatformTenant: h.platformTenantID != "" &&
			string(result.TenantId) == h.platformTenantID,
	})
}
