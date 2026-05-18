package commands

import (
	"context"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"
	"kali-auth-context/internal/ports"
)

type CreateUserDto struct {
	TenantId             identity.TenantId
	IdentificationNumber string
	Username             string
	Email                string
	Password             string
}

// Caso de uso para registrar un nuevo usuario en el sistema.
// Se encarga de validar la contraseña, generar el ID y guardar el usuario con la contraseña hasheada.
type CreateUserCommand struct {
	repo     ports.ICreateUserCommandRepository
	provider ports.IUUIDProvider
	hasher   ports.IPasswordHasher
}

func NewCreateUserCommand(repo ports.ICreateUserCommandRepository, provider ports.IUUIDProvider, hasher ports.IPasswordHasher) *CreateUserCommand {
	return &CreateUserCommand{repo: repo, provider: provider, hasher: hasher}
}

// Ejecuta el registro de usuario:
// 1. Valida la fortaleza de la contraseña.
// 2. Genera un nuevo ID de usuario.
// 3. Hashea la contraseña.
// 4. Crea el usuario en el repositorio.
func (c *CreateUserCommand) Execute(ctx context.Context, user *CreateUserDto) error {
       if err := policies.ValidatePasswordStrength(user.Password); err != nil {
	       // Si la contraseña no cumple la política, retorna error.
	       return err
       }

       id := identity.UserId(c.provider.Generate())

       hashedPassword, err := c.hasher.Hash(user.Password)
       if err != nil {
	       // Si ocurre error al hashear la contraseña, retorna error.
	       return err
       }

       domainUser, err := identity.NewUser(
	       id,
	       user.TenantId,
	       user.IdentificationNumber,
	       user.Username,
	       user.Email,
	       hashedPassword,
       )
       if err != nil {
	       // Si ocurre error al crear el usuario de dominio, retorna error.
	       return err
       }

       // Guarda el usuario en el repositorio.
       return c.repo.Create(ctx, domainUser)
}
