package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/EduGoGroup/edugo-shared/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	defaultMaxRetries  = 3
	defaultBaseBackoff = 100 * time.Millisecond
)

// Client implementa storage.Client para AWS S3.
// Incluye retry con backoff exponencial. NO incluye validacion de dominio
// (tipo de archivo, tamano maximo, etc.) — eso es responsabilidad del consumidor.
type Client struct {
	s3Client    *s3.Client
	bucket      string
	maxRetries  int
	baseBackoff time.Duration
}

// ClientOption configura el Client.
type ClientOption func(*Client)

// WithMaxRetries configura el numero maximo de reintentos.
func WithMaxRetries(n int) ClientOption {
	return func(c *Client) { c.maxRetries = n }
}

// WithBaseBackoff configura el backoff base para reintentos.
func WithBaseBackoff(d time.Duration) ClientOption {
	return func(c *Client) { c.baseBackoff = d }
}

// NewClient crea un storage.Client para S3.
// Recibe un *s3.Client ya configurado (via bootstrap/s3 o directamente).
func NewClient(s3Client *s3.Client, bucket string, opts ...ClientOption) *Client {
	c := &Client{
		s3Client:    s3Client,
		bucket:      bucket,
		maxRetries:  defaultMaxRetries,
		baseBackoff: defaultBaseBackoff,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	var lastErr error
	maxAttempts := 1 + c.maxRetries
	for attempt := range maxAttempts {
		if attempt > 0 {
			delay := c.baseBackoff * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("storage/s3: download %s: %w", key, ctx.Err())
			case <-time.After(delay):
			}
		}

		output, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(key),
		})
		if err == nil {
			return output.Body, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("storage/s3: download %s after %d retries: %w", key, c.maxRetries, lastErr)
}

// Upload sube contenido a S3. No hace retry porque io.Reader se consume
// en el primer intento y no es replayable de forma segura.
func (c *Client) Upload(ctx context.Context, key string, content io.Reader) error {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   content,
	})
	if err != nil {
		return fmt.Errorf("storage/s3: upload %s: %w", key, err)
	}
	return nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("storage/s3: delete %s: %w", key, err)
	}
	return nil
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Distinguir NotFound/NoSuchKey de errores reales (permisos, red, etc.)
		var nf *types.NotFound
		var nsk *types.NoSuchKey
		if errors.As(err, &nf) || errors.As(err, &nsk) {
			return false, nil
		}
		return false, fmt.Errorf("storage/s3: exists %s: %w", key, err)
	}
	return true, nil
}

func (c *Client) GetMetadata(ctx context.Context, key string) (*storage.FileMetadata, error) {
	output, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("storage/s3: metadata %s: %w", key, err)
	}

	var lastMod time.Time
	if output.LastModified != nil {
		lastMod = *output.LastModified
	}

	return &storage.FileMetadata{
		Key:          key,
		Size:         aws.ToInt64(output.ContentLength),
		ContentType:  aws.ToString(output.ContentType),
		LastModified: lastMod,
		ETag:         aws.ToString(output.ETag),
	}, nil
}

// Verificacion en compilacion.
var _ storage.Client = (*Client)(nil)
