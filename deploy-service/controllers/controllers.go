package controllers

import (
	"context"
	"github.com/Sumitk99/vercel/deploy-service/constants"
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

		err = srv.DownloadR2Folder(result[1])
		if err != nil {
			log.Printf("Error downloading folder: %v", err)
			continue
		}
	}
}
