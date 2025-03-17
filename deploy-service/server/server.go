package server

import (
	"context"
	"fmt"
	"github.com/Sumitk99/vercel/deploy-service/constants"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	client := s3.NewFromConfig(R2Config)

	return client, nil
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

func (srv *Server) DownloadR2Folder(ProjectID string) error {
	ProjectID += "/"

	objectList, err := srv.R2Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(constants.Bucket),
		Prefix: aws.String(ProjectID),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects: %v", err)
	}

	if len(objectList.Contents) == 0 {
		fmt.Println("Folder is empty or does not exist.")
		return nil
	}
	curPath, _ := os.Getwd()
	path := filepath.Join(curPath, constants.RepoPath, ProjectID)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create local directory: %v", err)
	}
	for _, obj := range objectList.Contents {
		relativePath := strings.TrimPrefix(*obj.Key, ProjectID)
		localFilePath := filepath.Join(curPath, relativePath)
		log.Println(relativePath, localFilePath)
		if err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create subdirectories: %v", err)
		}

		err := DownloadFile(srv.R2Client, *obj.Key, localFilePath)
		if err != nil {
			log.Printf("Failed to download %s: %v", *obj.Key, err)
		} else {
			fmt.Printf("Downloaded: %s -> %s\n", *obj.Key, localFilePath)
		}
	}

	fmt.Println("Folder downloaded successfully.")
	return nil
}

func DownloadFile(R2Client *s3.Client, key, localFilePath string) error {
	ctx := context.TODO()

	resp, err := R2Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(constants.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to download object: %v", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
