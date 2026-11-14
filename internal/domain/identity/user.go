package identity

import "errors"

type UserId string

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	// ErrEmailAlreadyExists indica que ya existe un usuario con ese email
	// dentro del mismo tenant (unicidad de email por tenant, no global).
	ErrEmailAlreadyExists = errors.New("a user with this email already exists for this tenant")
	// ErrUserHasRoles: no se elimina un usuario con roles asignados (sin cascada).
	ErrUserHasRoles = errors.New("user has roles assigned")
	// ErrCannotManageUser: editar o eliminar a otro usuario exige cubrir sus permisos.
	ErrCannotManageUser = errors.New("cannot manage a user with privileges you do not have")
	// ErrSelfDeletion: desde la administración de usuarios nadie se elimina a sí
	// mismo; otro administrador autorizado puede hacerlo.
	ErrSelfDeletion = errors.New("you cannot delete your own user")
)

type User struct {
	Id                   UserId
	TenantId             TenantId
	IdentificationNumber string
	Username             string
	Email                string
	Password             string
}

func NewUser(
	id UserId,
	tenantId TenantId,
	identificationNumber string,
	username string,
	email string,
	password string,
) (*User, error) {
	if tenantId == "" {
		return nil, ErrTenantRequired
	}

	if email == "" {
		return nil, ErrEmailRequired
	}

	if password == "" {
		return nil, ErrPasswordRequired
	}

	u := &User{
		Id:                   id,
		TenantId:             tenantId,
		IdentificationNumber: identificationNumber,
		Username:             username,
		Email:                email,
		Password:             password,
	}

	return u, nil
}

// NewTenantScopedUser fuerza la pertenencia del usuario a un tenant desde el dominio.
func NewTenantScopedUser(
	tenantId TenantId,
	id UserId,
	identificationNumber string,
	username string,
	email string,
	password string,
) (*User, error) {
	return NewUser(id, tenantId, identificationNumber, username, email, password)
}
