package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sumitk99/vercel/upload-service/constants"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"path/filepath"
)

type Server struct {
	R2Client    *s3.Client
	RedisClient *redis.Client
}

func ConnectToR2(AccessKeyID, SecretAccessKey, Endpoint string) (*s3.Client, error) {
	log.Println("Connecting to R2 : ", Endpoint, AccessKeyID, SecretAccessKey)
	R2Config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			AccessKeyID,     // Replace with your Cloudflare R2 Access Key ID
			SecretAccessKey, // Replace with your Cloudflare R2 Secret Access Key
			"",
		)),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: Endpoint}, nil
			},
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	R2Client := s3.NewFromConfig(R2Config)

	return R2Client, nil
}

func ConnectToRedis(Address string) (*redis.Client, error) {
	log.Println("connecting to redis")
	client := redis.NewClient(&redis.Options{
		Addr:     Address,
		Password: "",
		DB:       0,
	})
	res, err := client.Ping(context.Background()).Result()
	log.Println(res, err)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func UploadToR2(R2Client *s3.Client, baseDir string, Files []string) error {
	for _, file := range Files {
		newFile, err := os.Open(file)
		if err != nil {
			return errors.New("failed to open file")
		}
		defer newFile.Close()
		objectKey, err := filepath.Rel(baseDir, file)
		log.Println("Uploading file: ", objectKey)
		_, err = R2Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(constants.Bucket),
			Key:    aws.String(objectKey),
			Body:   newFile,
		})
		if err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}

	}
	return nil
}

func PushToRedis(RedisClient *redis.Client, projectId string) error {
	res := RedisClient.LPush(context.Background(), constants.BuildKey, projectId)
	err := res.Err()
	if err != nil {
		log.Println("Error pushing to Redis: ", err)
		return err
	}
	return nil
}
