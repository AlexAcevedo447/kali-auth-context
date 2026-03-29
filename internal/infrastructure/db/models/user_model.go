package models

import (
	"kali-auth-context/internal/domain/identity"
	"time"
)

type UserModel struct {
	Id                   identity.UserId `bun:",pk"`
	IdentificationNumber string
	Username             string
	Email                string
	Password             string
	CreatedAt            time.Time
}