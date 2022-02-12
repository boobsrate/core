package minio

import (
	"bytes"
	"context"
	"fmt"

	"github.com/boobsrate/core/internal/domain"
	"github.com/minio/minio-go/v7"
)

const (
	titsBucketName      = "tits"
	titsContentTypeJpeg = "image/jpeg"
)

type Storage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(client *minio.Client, bucketName string) *Storage {
	return &Storage{
		client:     client,
		bucketName: bucketName,
	}
}

func (t *Storage) CreateImageFromFile(ctx context.Context, imageName string, filePath string) error {
	_, err := t.client.FPutObject(
		ctx, titsBucketName, imageName, filePath, minio.PutObjectOptions{ContentType: titsContentTypeJpeg},
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

func (t *Storage) AssembleFileURL(tits *domain.Tits) {
	tits.URL = fmt.Sprintf("%s/%s/%s.webp", t.client.EndpointURL(), t.bucketName, tits.ID)
}
