//go:build wireinject

package wire

import (
	"github.com/google/wire"
	controllers "hieupc05.github/backend-server/internal/controller"
	uploadimage "hieupc05.github/backend-server/pkg/uploadImage"
)

func InitUploadRouterHandler(secretKey string) (*controllers.UploadController, error) {
	wire.Build(
		uploadimage.NewImgbbUpload,
		controllers.NewUploadController,
	)
	return new(controllers.UploadController), nil
}
