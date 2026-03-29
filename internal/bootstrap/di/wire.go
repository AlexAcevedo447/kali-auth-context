package bootstrap

import "github.com/google/wire"

var IdentitySet = wire.NewSet(
	SecuritySet,
)

func InitializeContainer() {
	wire.Build(
		IdentitySet,
	)
}
