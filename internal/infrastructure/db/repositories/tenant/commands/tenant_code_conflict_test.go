package db

import (
	"errors"
	"fmt"
	"testing"

	"kali-auth-context/internal/domain/identity"
)

// fakePgError imita pgdriver.Error (cuyo mapa de campos no es exportado):
// expone el SQLSTATE mediante Field('C').
type fakePgError struct{ code string }

func (e fakePgError) Error() string { return "pg error " + e.code }
func (e fakePgError) Field(k byte) string {
	if k == 'C' {
		return e.code
	}
	return ""
}

func TestTranslateTenantCodeConflict(t *testing.T) {
	other := errors.New("connection reset")
	cases := []struct {
		name string
		in   error
		want error
	}{
		{"nil stays nil", nil, nil},
		{"unique violation becomes ErrTenantCodeTaken", fakePgError{code: "23505"}, identity.ErrTenantCodeTaken},
		{"wrapped unique violation", fmt.Errorf("update: %w", fakePgError{code: "23505"}), identity.ErrTenantCodeTaken},
		{"other pg error passes through", fakePgError{code: "23503"}, fakePgError{code: "23503"}},
		{"non pg error passes through", other, other},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := translateTenantCodeConflict(tc.in)
			if got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
