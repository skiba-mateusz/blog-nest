package s3uploader

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/skiba-mateusz/blog-nest/config"
)

type S3Uploader struct {
	s3Client s3.Client
}

func New(region string) (*S3Uploader, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(region))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &S3Uploader{
		s3Client: *s3Client,
	}, nil
}

func (s *S3Uploader) PutObject(file io.Reader, filename, directory string) (string, error) {
	uniqueFilename := uuid.New().String() + filepath.Ext(filename)
	s3Key := filepath.Join(directory, uniqueFilename)

	_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(config.Envs.S3BucketName),
		Key:    aws.String(s3Key),
		Body: file,
 	})
	if err != nil {
		return "", err
	}

	return s3Key, nil
}