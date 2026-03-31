package providers

import (
	"kali-auth-context/internal/ports"

	"github.com/google/uuid"
)

type UUIDProvider struct{}

var _ ports.IUUIDProvider = (*UUIDProvider)(nil)

func NewUUIDProvider() *UUIDProvider {
	return &UUIDProvider{}
}

func (p *UUIDProvider) Generate() string {
	return uuid.New().String()
}
