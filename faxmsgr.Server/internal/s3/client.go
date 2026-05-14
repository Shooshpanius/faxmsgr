// Пакет s3 предоставляет клиент для S3-совместимого хранилища (MinIO).
package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client обёртка над MinIO-клиентом.
type Client struct {
	mc     *minio.Client
	bucket string
}

// NewClient создаёт S3-клиент для MinIO/S3.
// endpoint — адрес сервера (например, "localhost:9000"),
// accessKey / secretKey — учётные данные,
// useSSL — использовать HTTPS.
func NewClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*Client, error) {
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}
	return &Client{mc: mc, bucket: bucket}, nil
}

// Upload загружает файл в хранилище и возвращает URL объекта.
func (c *Client) Upload(ctx context.Context, objectName, contentType string, reader io.Reader, size int64) (string, error) {
	_, err := c.mc.PutObject(ctx, c.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload object: %w", err)
	}
	return fmt.Sprintf("/%s/%s", c.bucket, objectName), nil
}
