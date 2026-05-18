package commands

import (
	authcommands "kali-auth-context/internal/application/auth/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"
	"kali-auth-context/internal/infrastructure/security"

	"github.com/gofiber/fiber/v2"
)

// Handler encargado de manejar el inicio de sesión (login) de usuarios.
// Recibe las credenciales desde la petición HTTP y ejecuta el caso de uso de autenticación.
type LoginHandler struct {
	command     *authcommands.LoginCommand
	tokenIssuer *security.AccessTokenIssuer
}

func NewLoginHandler(command *authcommands.LoginCommand, tokenIssuer *security.AccessTokenIssuer) *LoginHandler {
	return &LoginHandler{command: command, tokenIssuer: tokenIssuer}
}

type loginRequest struct {
	TenantId string `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string              `json:"access_token"`
	TokenType   string              `json:"token_type"`
	ExpiresIn   int64               `json:"expires_in"`
	TenantId    string              `json:"tenant_id"`
	UserId      string              `json:"user_id"`
	Email       string              `json:"email"`
	NeedsRehash bool                `json:"needs_rehash"`
	Roles       []string            `json:"roles"`
	Permissions []map[string]string `json:"permissions"`
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
	       TenantId: identity.TenantId(req.TenantId),
	       Email:    req.Email,
	       Password: req.Password,
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
	       UserId:      string(result.UserId),
	       Email:       result.Email,
	       NeedsRehash: result.NeedsRehash,
	       Roles:       roleNames,
	       Permissions: perms,
       })
}
