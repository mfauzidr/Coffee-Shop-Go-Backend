package pkg

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type Cloudinary struct {
	CLD *cloudinary.Cloudinary
}

func NewCloudinaryUtil() *Cloudinary {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatal("Failed to initiate Cloudinary: %w", err)
	}

	return &Cloudinary{
		CLD: cld,
	}
}

func (c *Cloudinary) UploadFile(ctx *gin.Context, file interface{}, fileName string) (*uploader.UploadResult, error) {
	uploadResult, err := c.CLD.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileName,
	})
	if err != nil {
		return nil, err
	}
	return uploadResult, nil
}
