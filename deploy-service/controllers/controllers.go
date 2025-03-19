package controllers

import (
	"context"
	"encoding/json"
	"github.com/Sumitk99/vercel/deploy-service/builder"
	"github.com/Sumitk99/vercel/deploy-service/constants"
	"github.com/Sumitk99/vercel/deploy-service/helper"
	"github.com/Sumitk99/vercel/deploy-service/models"
	"github.com/Sumitk99/vercel/deploy-service/server"
	"log"

	"time"
)

func StartRedisListener(srv *server.Server) {
	for {
		result, err := srv.RedisClient.BLPop(context.Background(), 0*time.Second, constants.BuildKey).Result()
		if err != nil {
			log.Printf("Error fetching from queue: %v", err)
			continue
		}
		log.Println(result)
		object := &models.RedisObject{}
		err = json.Unmarshal([]byte(result[1]), object)
		projectPath, err := srv.DownloadR2Folder(object.ProjectId)
		if err != nil {
			log.Printf("Error downloading folder: %v", err)
			continue
		}
		log.Println(object, object.ProjectId)
		buildPath, err := builder.BuildAngularProject(*projectPath)
		files := helper.GetAllFiles(*buildPath)
		log.Println(files)
		_ = server.UploadBuildToR2(srv.R2Client, *buildPath, object.ProjectId, files)
	}
}
