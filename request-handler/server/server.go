package server

import (
	"context"
	"fmt"
	"github.com/Sumitk99/vercel/request-handler/constants"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
	"log"
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
			AccessKeyID,
			SecretAccessKey,
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

func DownloadFileFromR2(R2Client *s3.Client, key string) (*s3.GetObjectOutput, error) {

	resp, err := R2Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download object: %v", err)
	}
	return resp, nil
}
