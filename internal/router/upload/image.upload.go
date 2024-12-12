package upload

import (
	"log"

	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/wire"
)

type UploadImageRouter struct{}

func (r *UploadImageRouter) InitUploadImageRouter(Router *gin.RouterGroup) {

	uploadHandlerNonDependency, err := wire.InitUploadRouterHandler(global.Config.Imgbb.ApiKey)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	uploadImageRouter := Router.Group("/upload")
	{
		uploadImageRouter.POST("/image", uploadHandlerNonDependency.UploadImage)
	}
}
