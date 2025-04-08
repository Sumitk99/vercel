package server

import (
	"context"
	"errors"
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
	"sync"
	"time"
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

func (srv *Server) DownloadR2Folder(ProjectID string) (*string, error) {
	ProjectID += "/"

	objectList, err := srv.R2Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(constants.Bucket),
		Prefix: aws.String(ProjectID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	if len(objectList.Contents) == 0 {
		fmt.Println("Folder is empty or does not exist.")
		return nil, errors.New("folder is empty or does not exist")
	}
	curPath, _ := os.Getwd()
	path := filepath.Join(curPath, constants.RepoPath, ProjectID)
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create local directory: %v", err)
	}
	log.Println(path)
	start := time.Now()
	wg := &sync.WaitGroup{}
	for _, obj := range objectList.Contents {
		log.Println(*obj.Key)
		wg.Add(1)

		go func(WaitGroup *sync.WaitGroup) {
			relativePath := strings.TrimPrefix(*obj.Key, ProjectID)
			localFilePath := filepath.Join(path, relativePath)
			log.Println(relativePath, localFilePath)
			if err = os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm); err != nil {
				log.Println("failed to create subdirectories: %v", err)
			}

			err = DownloadFileFromR2(srv.R2Client, *obj.Key, localFilePath)
			if err != nil {
				log.Printf("Failed to download %s: %v", *obj.Key, err)
			} else {
				fmt.Printf("Downloaded: %s -> %s\n", *obj.Key, localFilePath)
			}
			WaitGroup.Done()
		}(wg)

	}

	wg.Wait()
	time.Since(start)
	fmt.Println("Folder downloaded successfully in ", time.Since(start), "secs")
	return &path, nil
}

func DownloadFileFromR2(R2Client *s3.Client, key, localFilePath string) error {
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

func UploadBuildToR2(R2Client *s3.Client, baseDir, projectId string, Files []string) error {
	start := time.Now()
	log.Println("Base Directory : ", baseDir)
	wg := &sync.WaitGroup{}
	for _, file := range Files {
		wg.Add(1)
		go func(WaitGroup *sync.WaitGroup) {
			newFile, _ := os.Open(file)
			defer newFile.Close()
			OriginalObjectKey, _ := filepath.Rel(baseDir, file)
			outputPath := filepath.Join(constants.OutputPath, projectId, OriginalObjectKey)
			log.Println("Uploading file: ", outputPath)
			_, err := R2Client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: aws.String(constants.Bucket),
				Key:    aws.String(outputPath),
				Body:   newFile,
			})
			if err != nil {
				log.Println(err)
				return
			}
			WaitGroup.Done()
		}(wg)
	}
	wg.Wait()
	log.Printf("Uploading %v files took %s secs\n", len(Files), time.Since(start))
	return nil
}
