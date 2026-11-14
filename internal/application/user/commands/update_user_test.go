package commands

import (
	"context"
	"errors"
	"testing"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"
)

type captureUpdateRepo struct{ got *identity.User }

func (r *captureUpdateRepo) Update(_ context.Context, u *identity.User) error {
	r.got = u
	return nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(p string) (string, error) { return "hash:" + p, nil }
func (fakeHasher) Compare(string, string) error  { return nil }
func (fakeHasher) NeedsRehash(string) bool       { return false }

func TestUpdateUser_EmptyPasswordKeepsCurrent(t *testing.T) {
	repo := &captureUpdateRepo{}
	err := NewUpdateUserCommand(repo, fakeHasher{}).Execute(context.Background(), &UpdateUserDto{
		Id: "u1", TenantId: "t1", Username: "ana", Email: " ana@x.com ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.got.Password != "" {
		t.Fatalf("empty password must not be sent to the repository, got %q", repo.got.Password)
	}
	if repo.got.Email != "ana@x.com" {
		t.Fatalf("email must be trimmed, got %q", repo.got.Email)
	}
}

func TestUpdateUser_NewPasswordIsValidatedAndHashed(t *testing.T) {
	repo := &captureUpdateRepo{}
	cmd := NewUpdateUserCommand(repo, fakeHasher{})
	if err := cmd.Execute(context.Background(), &UpdateUserDto{Id: "u1", TenantId: "t1", Email: "a@x.com", Password: "corta"}); !errors.Is(err, policies.ErrPasswordTooShort) {
		t.Fatalf("expected password policy error, got %v", err)
	}
	if err := cmd.Execute(context.Background(), &UpdateUserDto{Id: "u1", TenantId: "t1", Email: "a@x.com", Password: "Nueva#Clave2026"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.got.Password != "hash:Nueva#Clave2026" {
		t.Fatalf("password must be hashed, got %q", repo.got.Password)
	}
}

func TestUpdateUser_EmailRequired(t *testing.T) {
	err := NewUpdateUserCommand(&captureUpdateRepo{}, fakeHasher{}).Execute(context.Background(), &UpdateUserDto{Id: "u1", TenantId: "t1", Email: "  "})
	if !errors.Is(err, identity.ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}
