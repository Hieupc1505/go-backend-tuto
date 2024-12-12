package uploadimage

import (
	"encoding/base64"
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

const ImgurAPI = "https://api.imgur.com/3/image"

type ImgurUpload struct {
	ClientID string
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

func NewImgurUpload(clientID string) IUploadImage {
	return &ImgurUpload{
		ClientID: clientID,
	}
}

func (i *ImgurUpload) Upload(image string) (UploadResult, error) {

	imgBytes, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return UploadResult{}, err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Client-ID "+i.ClientID).
		SetBody(imgBytes).
		Post(ImgurAPI)
	if err != nil {
		return UploadResult{}, err
	}

	//parse response
	var imgurResponse ImgurResponse
	if err := json.Unmarshal(resp.Body(), &imgurResponse); err != nil {
		return UploadResult{}, err
	}

	return UploadResult{Url: imgurResponse.Data.Link, Thumb: ""}, nil

}
