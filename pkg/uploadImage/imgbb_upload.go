package uploadimage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

const ImgbbAPI = "https://api.imgbb.com/1/upload"

type ImgbbUpload struct {
	ApiKey string
}

// ImgbbResponse defines the structure of the response from Imgbb
type ImgbbResponse struct {
	Data    ImgbbData `json:"data"`
	Success bool      `json:"success"`
	Status  int       `json:"status"`
}

type ImgbbData struct {
	Url   string    `json:"url"`
	Thumb ThumbData `json:"thumb"`
}

type ThumbData struct {
	Url string `json:"url"`
}

func NewImgbbUpload(apikey string) IUploadImage {
	return &ImgbbUpload{
		ApiKey: apikey,
	}
}

func (i *ImgbbUpload) Upload(image string) (UploadResult, error) {

	imageBytes, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return UploadResult{}, err
	}

	// Upload to Imgbb
	client := resty.New()
	resp, err := client.R().
		SetFileReader("image", "image.jpg", bytes.NewReader(imageBytes)).
		SetQueryParam("key", i.ApiKey).
		Post(ImgbbAPI)

	if err != nil {
		return UploadResult{}, err
	}

	//parse response
	var ImgbbResponse ImgbbResponse
	if err := json.Unmarshal(resp.Body(), &ImgbbResponse); err != nil {
		return UploadResult{}, err
	}

	return UploadResult{Url: ImgbbResponse.Data.Url, Thumb: ImgbbResponse.Data.Thumb.Url}, nil

}
