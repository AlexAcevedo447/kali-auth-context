package bootstrap

import (
	"kali-auth-context/internal/infrastructure/security"
	"kali-auth-context/internal/ports"

	"github.com/google/wire"
)

var SecuritySet = wire.NewSet(
	security.NewBcryptHasher,
	wire.Bind(new(ports.IPasswordHasher), new(*security.BcryptHasher)),
)

var CommandSet = wire.NewSet(
)
