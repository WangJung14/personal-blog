package utils

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(filePath string) (string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{
		Folder: "blog_uploads",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func UploadFileToCloudinary(file interface{}) (string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "blog_uploads",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
