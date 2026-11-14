package shared

import "kali-auth-context/internal/domain/identity"

// UserResponse es la forma pública de un usuario: NUNCA incluye la contraseña
// ni su hash. Conserva las claves PascalCase que ya consumen los clientes.
type UserResponse struct {
	Id                   identity.UserId   `json:"Id"`
	TenantId             identity.TenantId `json:"TenantId"`
	IdentificationNumber string            `json:"IdentificationNumber"`
	Username             string            `json:"Username"`
	Email                string            `json:"Email"`
}

func ToUserResponse(u *identity.User) UserResponse {
	return UserResponse{
		Id:                   u.Id,
		TenantId:             u.TenantId,
		IdentificationNumber: u.IdentificationNumber,
		Username:             u.Username,
		Email:                u.Email,
	}
}

func ToUserResponses[T interface{ ~*identity.User }](users []T) []UserResponse {
	out := make([]UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, ToUserResponse((*identity.User)(u)))
	}
	return out
}
