//go:build wireinject

package wire

import (
	"github.com/google/wire"
	controllers "hieupc05.github/backend-server/internal/controller"
	repos "hieupc05.github/backend-server/internal/repo"
	"hieupc05.github/backend-server/internal/services"
	"hieupc05.github/backend-server/internal/utils/token"
)

func InitUserRouterHandler(secretKey string, tokenMaker token.Maker) (*controllers.UserController, error) {
	wire.Build(
		repos.NewUserRepository,
		repos.NewUserAuthRepository,
		services.NewUserServices,
		controllers.NewUserController,
	)
	return new(controllers.UserController), nil
}
