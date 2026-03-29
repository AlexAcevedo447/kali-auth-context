package identity

import "errors"

type UserId string

type User struct {
	Id                   UserId
	IdentificationNumber string
	Username             string
	Email                string
	Password             string
}

func NewUser(
	Id UserId,
	identificationNumber string,
	username string,
	email string,
	password string,
) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	u := &User{
		Id:                   Id,
		IdentificationNumber: identificationNumber,
		Username:             username,
		Email:                email,
		Password:             password,
	}

	return u, nil
}
