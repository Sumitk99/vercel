package controllers

import (
	"fmt"
	"github.com/Sumitk99/vercel/upload-service/constants"
	"github.com/Sumitk99/vercel/upload-service/helper"
	"os"
	"path/filepath"

	"github.com/Sumitk99/vercel/upload-service/models"
	"github.com/Sumitk99/vercel/upload-service/server"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
	"os/exec"
)

func Controller(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Req
		err := c.BindJSON(&req)
		if err != nil {
			log.Println(err.Error())
		}
		projectId := ksuid.New().String()[:5]

		url := req.RepoUrl
		go func() {
			err = helper.CloneRepo(url, projectId)
			if err != nil {
				log.Println(err)
				return
			}
			rootDirectory, _ := os.Getwd()
			baseDir := filepath.Join(rootDirectory, constants.RepoPath)
			dir := filepath.Join(baseDir, projectId)
			files := helper.GetAllFiles(dir)
			err = server.UploadToR2(srv.R2Client, baseDir, files)
			if err != nil {
				log.Println(err)
				return
			}
			err = exec.Command("rm", "-rf", fmt.Sprintf("%s/%s", constants.RepoPath, projectId)).Run()
			log.Println("Pushing to Redis Queue")
			err = server.PushToRedis(srv.RedisClient, projectId, req.Framework)
			log.Println("Pushed to Redis Queue")
		}()

		c.JSON(http.StatusOK, gin.H{
			"url":       req.RepoUrl,
			"projectId": projectId,
			"status":    "Clone in Progress",
		})
	}
}

func FetchStatus(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectId := c.Param("projectId")
		//status, err := server.FetchStatus(srv.RedisClient, projectId)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//}
		c.JSON(http.StatusOK, gin.H{
			"projectId": projectId,
			//"status":    status,
		})
	}
}
