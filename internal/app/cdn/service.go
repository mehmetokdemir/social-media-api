package cdn

import (
	"context"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"mime/multipart"
)

type ICdnService interface {
	UploadImage(file *multipart.FileHeader) (string, error)
}

type cdnService struct {
	cloudinaryClient *cloudinary.Cloudinary
}

func NewCdnService(cloudinaryClient *cloudinary.Cloudinary) ICdnService {
	if cloudinaryClient == nil {
		return nil
	}

	return &cdnService{
		cloudinaryClient: cloudinaryClient,
	}
}

func (s *cdnService) UploadImage(file *multipart.FileHeader) (string, error) {
	fileToUpload, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileToUpload.Close()

	ctx := context.Background()
	uploadResult, err := s.cloudinaryClient.Upload.Upload(ctx, fileToUpload, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.URL, nil
}
