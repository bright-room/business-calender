package awsx

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/xerrors"
)

type S3Client struct {
	client     s3.Client
	uploader   manager.Uploader
	downloader manager.Downloader
	bucket     string
}

func (r *S3Client) Upload(reader io.Reader, key string, ctx context.Context) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	if _, err := r.uploader.Upload(ctx, input); err != nil {
		return xerrors.Errorf("failed to s3 upload[%s/%s]: %w", r.bucket, key, err)
	}

	return nil
}

func (r *S3Client) Download(key string, ctx context.Context) (io.Reader, error) {
	var file *os.File
	var err error

	file, err = os.CreateTemp("", filepath.Base(key))
	if err != nil {
		return nil, xerrors.Errorf("failed to create file[%s]: %w", filepath.Base(key), err)
	}

	defer func() {
		if err2 := file.Close(); err2 != nil {
			err = xerrors.Errorf("failed to close file[%s]: %w", filepath.Base(key), err)
		}
	}()

	input := &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	}

	if _, err := r.downloader.Download(ctx, file, input); err != nil {
		return nil, xerrors.Errorf("failed to s3 download[%s/%s]: %w", r.bucket, key, err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerrors.Errorf("failed to read file[%s]: %w", filepath.Base(key), err)
	}

	return bytes.NewReader(data), nil
}

func NewS3Client(bucket string, awsConfig aws.Config, optFns ...func(options *s3.Options)) *S3Client {
	client := s3.NewFromConfig(awsConfig, optFns...)

	uploader := manager.NewUploader(client)
	uploader.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(5 * 1024 * 1024)

	downloader := manager.NewDownloader(client)
	downloader.BufferProvider = manager.NewPooledBufferedWriterReadFromProvider(5 * 1024 * 1024)

	return &S3Client{
		client:     *client,
		uploader:   *uploader,
		downloader: *downloader,
		bucket:     bucket,
	}
}
