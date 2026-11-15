package policies

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

// MinPasswordLength se mide en caracteres Unicode (code points), no en bytes:
// "Ñandú#Clav1" tiene 11 caracteres aunque ocupe 13 bytes en UTF-8.
// No hay máximo propio: las contraseñas nuevas se guardan con Argon2id, que usa
// la contraseña completa. bcrypt (límite efectivo de 72 bytes) solo verifica
// hashes heredados "$2…" (el login los señala con needs_rehash; no se migran solos).
const MinPasswordLength = 12

var (
	ErrPasswordTooShort      = errors.New("password must have at least 12 characters")
	ErrPasswordMissingUpper  = errors.New("password must include at least one uppercase letter")
	ErrPasswordMissingLower  = errors.New("password must include at least one lowercase letter")
	ErrPasswordMissingDigit  = errors.New("password must include at least one digit")
	ErrPasswordMissingSymbol = errors.New("password must include at least one symbol")
	ErrPasswordHasWhitespace = errors.New("password must not contain whitespace")
)

func ValidatePasswordStrength(password string) error {
	if utf8.RuneCountInString(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSymbol := false

	for _, c := range password {
		if unicode.IsSpace(c) {
			return ErrPasswordHasWhitespace
		}
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		if unicode.IsLower(c) {
			hasLower = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}
		if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			hasSymbol = true
		}
	}

	if !hasUpper {
		return ErrPasswordMissingUpper
	}
	if !hasLower {
		return ErrPasswordMissingLower
	}
	if !hasDigit {
		return ErrPasswordMissingDigit
	}
	if !hasSymbol {
		return ErrPasswordMissingSymbol
	}

	return nil
}
