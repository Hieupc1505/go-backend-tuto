package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/global"
	uploadimage "hieupc05.github/backend-server/pkg/uploadImage"
)

type UploadController struct {
	uploadService uploadimage.IUploadImage
}

type UploadImageRequest struct {
	Image string `json:"image"`
}

// ImgurResponse defines the structure of the response from Imgur
type ImgurResponse struct {
	Data    ImgurData `json:"data"`
	Success bool      `json:"success"`
	Status  int       `json:"status"`
}

type ImgurData struct {
	Link string `json:"link"`
}

func NewUploadController(upload uploadimage.IUploadImage) *UploadController {
	return &UploadController{
		uploadService: upload,
	}
}

func (uc *UploadController) UploadImage(ctx *gin.Context) {
	var payload UploadImageRequest
	clientID := global.Config.Imgur.ClientID
	if clientID == "" {
		log.Fatal("IMGUR_CLIENT_ID is not set in .env file")
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Image == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	data, err := uc.uploadService.Upload(payload.Image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)

}
