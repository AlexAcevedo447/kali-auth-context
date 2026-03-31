package policies

import (
	"errors"
	"unicode"
)

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
	if len(password) < MinPasswordLength {
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
