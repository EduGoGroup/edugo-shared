package s3

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-shared/bootstrap"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Factory implementa la creacion de clientes S3.
type Factory struct{}

// NewFactory crea una nueva Factory de S3.
func NewFactory() *Factory {
	return &Factory{}
}

// CreateClient crea un cliente S3 configurado.
// Soporta endpoints custom (LocalStack) via S3Config.Endpoint.
func (f *Factory) CreateClient(ctx context.Context, cfg bootstrap.S3Config) (*s3.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"",
	)

	awsOpts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(creds),
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsOpts...)
	if err != nil {
		return nil, fmt.Errorf("bootstrap/s3: load AWS config: %w", err)
	}

	s3Opts := []func(*s3.Options){}
	if cfg.Endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = cfg.ForcePathStyle
		})
	}

	client := s3.NewFromConfig(awsCfg, s3Opts...)

	if cfg.Bucket != "" {
		if err := f.ValidateBucket(ctx, client, cfg.Bucket); err != nil {
			return nil, fmt.Errorf("bootstrap/s3: bucket validation: %w", err)
		}
	}

	return client, nil
}

// CreatePresignClient crea un cliente para URLs pre-firmadas.
func (f *Factory) CreatePresignClient(client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(client)
}

// ValidateBucket verifica que el bucket existe y es accesible.
func (f *Factory) ValidateBucket(ctx context.Context, client *s3.Client, bucket string) error {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("bucket '%s' not accessible: %w", bucket, err)
	}
	return nil
}
