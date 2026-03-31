package identity

import "errors"

type UserId string

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
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
