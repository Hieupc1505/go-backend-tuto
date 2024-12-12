package services

import uploadimage "hieupc05.github/backend-server/pkg/uploadImage"

type UploadService struct {
	Upload uploadimage.IUploadImage
}

func NewUploadService(upload uploadimage.IUploadImage) *UploadService {
	return &UploadService{
		Upload: upload,
	}
}

func (s *UploadService) UploadImage(image string) (uploadimage.UploadResult, error) {
	return s.Upload.Upload(image)
}
