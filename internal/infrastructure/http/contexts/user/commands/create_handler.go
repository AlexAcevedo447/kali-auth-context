package commands

import (
	usercommands "kali-auth-context/internal/application/user/commands"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

// Handler encargado de registrar un nuevo usuario.
// Recibe los datos del usuario desde la petición HTTP y ejecuta el caso de uso de registro.
type CreateHandler struct {
	command *usercommands.CreateUserCommand
}

func NewCreateHandler(command *usercommands.CreateUserCommand) *CreateHandler {
	return &CreateHandler{command: command}
}

// Maneja la petición POST para registrar un usuario.
// Los datos esperados son tenant_id, identification_number, username, email y password en el cuerpo JSON.
// Si el registro es exitoso, retorna 201 Created. Si hay error de validación, retorna el error correspondiente.
func (h *CreateHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		TenantId             string `json:"tenant_id"`
		IdentificationNumber string `json:"identification_number"`
		Username             string `json:"username"`
		Email                string `json:"email"`
		Password             string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Ejecuta el caso de uso de registro de usuario con los datos recibidos.
	err := h.command.Execute(c.UserContext(), &usercommands.CreateUserDto{
		TenantId:             identity.TenantId(req.TenantId),
		IdentificationNumber: req.IdentificationNumber,
		Username:             req.Username,
		Email:                req.Email,
		Password:             req.Password,
	})
	if err != nil {
		// Si ocurre un error de validación o negocio, se retorna el error correspondiente.
		return shared.WriteError(c, err)
	}

	// Si el registro es exitoso, retorna 201 Created.
	return c.SendStatus(fiber.StatusCreated)
}
