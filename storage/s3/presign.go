package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-shared/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const defaultPresignExpiry = 15 * time.Minute

// PresignClient implementa storage.PresignClient para AWS S3.
type PresignClient struct {
	presigner *s3.PresignClient
	bucket    string
	expiry    time.Duration
}

// NewPresignClient crea un storage.PresignClient para S3.
// Recibe un *s3.Client ya configurado (via bootstrap/s3 o directamente).
// Si expiry es 0, usa 15 minutos por defecto.
func NewPresignClient(s3Client *s3.Client, bucket string, expiry time.Duration) *PresignClient {
	if expiry == 0 {
		expiry = defaultPresignExpiry
	}
	return &PresignClient{
		presigner: s3.NewPresignClient(s3Client),
		bucket:    bucket,
		expiry:    expiry,
	}
}

func (c *PresignClient) GenerateUploadURL(ctx context.Context, key string) (string, time.Time, error) {
	result, err := c.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(c.expiry))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("storage/s3: presign PUT %s: %w", key, err)
	}
	return result.URL, time.Now().Add(c.expiry), nil
}

func (c *PresignClient) GenerateDownloadURL(ctx context.Context, key string) (string, time.Time, error) {
	result, err := c.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(c.expiry))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("storage/s3: presign GET %s: %w", key, err)
	}
	return result.URL, time.Now().Add(c.expiry), nil
}

// Verificacion en compilacion.
var _ storage.PresignClient = (*PresignClient)(nil)
