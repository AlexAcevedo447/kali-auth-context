package mappers

import (
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
)

func ToDomainUser(userModel *models.UserModel) (*identity.User, error) {
	return identity.NewUser(
		identity.UserId(userModel.Id),
		identity.TenantId(userModel.TenantId),
		userModel.IdentificationNumber,
		userModel.Username,
		userModel.Email,
		userModel.Password,
	)
}

func ToUserModel(user *identity.User) *models.UserModel {
	return &models.UserModel{
		Id:                   identity.UserId(user.Id),
		TenantId:             identity.TenantId(user.TenantId),
		IdentificationNumber: user.IdentificationNumber,
		Username:             user.Username,
		Email:                user.Email,
		Password:             user.Password,
	}
}
