package minio

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

const (
	titsContentTypeJpeg = "image/jpeg"
)

type Storage struct {
	client     *minio.Client
	bucketName string
	publicURL  string
}

func NewMinioStorage(client *minio.Client, bucketName string, publicURL string) *Storage {
	return &Storage{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicURL,
	}
}

func (t *Storage) CreateImageFromFile(ctx context.Context, imageName string, filePath string) error {
	_, err := t.client.FPutObject(
		ctx, t.bucketName, imageName, filePath, minio.PutObjectOptions{ContentType: titsContentTypeJpeg},
	)
	if err != nil {
		return fmt.Errorf("upload image to minio: %v", err)
	}
	return nil
}

func (t *Storage) CreateImageFromBytes(ctx context.Context, imageName string, imageData []byte) error {

	reader := bytes.NewReader(imageData)

	_, err := t.client.PutObject(
		ctx, t.bucketName, imageName, reader, int64(len(imageData)), minio.PutObjectOptions{ContentType: titsContentTypeJpeg},
	)
	if err != nil {
		return fmt.Errorf("upload image to minio: %v", err)
	}
	return nil
}

func (t *Storage) DeleteImage(ctx context.Context, imageName string) error {
	err := t.client.RemoveObject(ctx, t.bucketName, imageName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return fmt.Errorf("delete image from minio: %v", err)
	}
	return nil
}

func (t *Storage) GetImageUrl(imageID string) string {
	return fmt.Sprintf("%s/%s/%s", t.publicURL, t.bucketName, imageID)
}
