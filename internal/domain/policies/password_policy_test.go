package policies

import (
	"errors"
	"testing"
)

// El mínimo de 12 se cuenta en caracteres Unicode (code points), no en bytes,
// igual que la interfaz ([...password].length).
func TestValidatePasswordStrength(t *testing.T) {
	cases := []struct {
		name     string
		password string
		want     error
	}{
		{"ASCII de 12 caracteres", "Clave#2026ab", nil},
		{"ASCII de 11 caracteres", "Clave#2026a", ErrPasswordTooShort},
		{"11 caracteres con ñ y tilde (13 bytes)", "Ñandú#Clav1", ErrPasswordTooShort},
		{"12 caracteres con ñ y tilde (14 bytes)", "Ñandú#Clave1", nil},
		{"11 code points con emoji (14 bytes)", "Clave#2026😀", ErrPasswordTooShort},
		{"12 code points con emoji (15 bytes)", "Clave#2026😀x", nil},
		{"sin mayúscula", "clave#2026abc", ErrPasswordMissingUpper},
		{"sin minúscula", "CLAVE#2026ABC", ErrPasswordMissingLower},
		{"sin dígito", "Clave#Segura!", ErrPasswordMissingDigit},
		{"sin símbolo", "ClaveSegura2026", ErrPasswordMissingSymbol},
		{"con espacio", "Clave #2026ab", ErrPasswordHasWhitespace},
		{"con espacio no separable", "Clave #2026ab", ErrPasswordHasWhitespace},
		{"mayúscula y minúscula no ASCII", "ÁÉÍÓÚñ#2026xx", nil},
		{"larga de más de 72 bytes", "Ñandú#Clave1" + "ññññññññññññññññññññññññññññññññññññññ", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ValidatePasswordStrength(tc.password); !errors.Is(got, tc.want) {
				t.Fatalf("ValidatePasswordStrength(%q) = %v, want %v", tc.password, got, tc.want)
			}
		})
	}
}
