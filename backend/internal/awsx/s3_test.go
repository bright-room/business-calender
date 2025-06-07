package awsx_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"net.bright-room.dev/calender-api/internal/awsx"
)

func TestS3Client_Upload(t *testing.T) {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("AWS_ENDPOINT")
	awsBucket := os.Getenv("AWS_BUCKET")

	credential := credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")
	awsConfig, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithBaseEndpoint(endpoint),
		config.WithCredentialsProvider(credential),
	)

	client := awsx.NewS3Client(awsBucket, awsConfig, func(options *s3.Options) { options.UsePathStyle = true })

	testData := "testtesttest"
	err := client.Upload(bytes.NewReader([]byte(testData)), "test.txt", context.Background())

	assert.NoError(t, err)
}

func TestS3Client_Download(t *testing.T) {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("AWS_ENDPOINT")
	awsBucket := os.Getenv("AWS_BUCKET")

	credential := credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")
	awsConfig, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithBaseEndpoint(endpoint),
		config.WithCredentialsProvider(credential),
	)

	client := awsx.NewS3Client(awsBucket, awsConfig, func(options *s3.Options) { options.UsePathStyle = true })

	reader, err := client.Download("download_test.txt", context.Background())
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	assert.NoError(t, err)

	actual := buf.String()
	assert.Equal(t, "Download test.", actual)
}
