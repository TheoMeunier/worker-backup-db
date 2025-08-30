package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

type ServiceS3Impl struct {
	s3Client *s3.Client
	bucket   string
}

func NewServiceS3WithR2() (*ServiceS3Impl, error) {
	_ = godotenv.Load()
	endpointURL := os.Getenv("S3_ENDPOINT_URL")

	cfg := aws.Config{
		Region: "auto",
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
			"",
		),
		DefaultsMode: aws.DefaultsModeStandard,
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpointURL)
		o.UsePathStyle = true
		o.Credentials = cfg.Credentials
	})

	return &ServiceS3Impl{
		s3Client: s3Client,
		bucket:   os.Getenv("S3_BUCKET"),
	}, nil
}

func (s *ServiceS3Impl) UploadToS3(data []byte, database string) error {
	currentTime := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("backup-%s-%s.sql.gz", database, currentTime)
	fullPath := fmt.Sprintf("%s/%s", database, filename)

	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fullPath),
		Body:   bytes.NewReader(data),
	}

	_, err := s.s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("error upload S3: %v", err)
	}
	return nil
}

func (s *ServiceS3Impl) DeleteFromS3(filePath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	}

	_, err := s.s3Client.DeleteObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("error remove file S3: %v", err)
	}

	return nil
}
