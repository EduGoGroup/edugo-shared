package bootstrap

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// =============================================================================
// S3 FACTORY IMPLEMENTATION
// =============================================================================

// DefaultS3Factory implementa S3Factory usando AWS SDK v2
type DefaultS3Factory struct{}

// NewDefaultS3Factory crea una nueva instancia de DefaultS3Factory
func NewDefaultS3Factory() *DefaultS3Factory {
	return &DefaultS3Factory{}
}

// CreateClient crea un cliente de S3
func (f *DefaultS3Factory) CreateClient(ctx context.Context, s3Config S3Config) (*s3.Client, error) {
	// Configurar credenciales estáticas
	creds := credentials.NewStaticCredentialsProvider(
		s3Config.AccessKeyID,
		s3Config.SecretAccessKey,
		"", // session token (vacío para credenciales estáticas)
	)

	// Cargar configuración AWS
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s3Config.Region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Crear cliente S3
	client := s3.NewFromConfig(cfg)

	// Validar que el bucket existe
	if err := f.ValidateBucket(ctx, client, s3Config.Bucket); err != nil {
		return nil, fmt.Errorf("bucket validation failed: %w", err)
	}

	return client, nil
}

// CreatePresignClient crea un cliente para URLs pre-firmadas
func (f *DefaultS3Factory) CreatePresignClient(client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(client)
}

// ValidateBucket verifica que el bucket existe y es accesible
func (f *DefaultS3Factory) ValidateBucket(ctx context.Context, client *s3.Client, bucket string) error {
	// Intentar head bucket
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("bucket '%s' not accessible: %w", bucket, err)
	}

	return nil
}

// Verificar que DefaultS3Factory implementa S3Factory
var _ S3Factory = (*DefaultS3Factory)(nil)
