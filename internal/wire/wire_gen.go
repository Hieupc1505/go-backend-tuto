// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"hieupc05.github/backend-server/internal/controller"
	"hieupc05.github/backend-server/internal/repo"
	"hieupc05.github/backend-server/internal/services"
	"hieupc05.github/backend-server/internal/utils/token"
	"hieupc05.github/backend-server/pkg/uploadImage"
)

// Injectors from upload.wire.go:

func InitUploadRouterHandler(secretKey string) (*controllers.UploadController, error) {
	iUploadImage := uploadimage.NewImgbbUpload(secretKey)
	uploadController := controllers.NewUploadController(iUploadImage)
	return uploadController, nil
}

// Injectors from user.wire.go:

func InitUserRouterHandler(secretKey string, tokenMaker token.Maker) (*controllers.UserController, error) {
	iUserRepository := repos.NewUserRepository()
	iUserAuthRepository := repos.NewUserAuthRepository()
	iUserServices := services.NewUserServices(iUserRepository, iUserAuthRepository, tokenMaker)
	userController := controllers.NewUserController(iUserServices)
	return userController, nil
}
