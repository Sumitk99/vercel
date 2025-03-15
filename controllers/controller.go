package controllers

import (
	"github.com/Sumitk99/vercel/constants"
	"github.com/Sumitk99/vercel/helper"
	"github.com/Sumitk99/vercel/models"
	"github.com/Sumitk99/vercel/server"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
			helper.CloneRepo(url, projectId)
			rootDirectory, _ := os.Getwd()
			baseDir := filepath.Join(rootDirectory, constants.RepoPath)
			dir := filepath.Join(baseDir, projectId)
			files := helper.GetAllFiles(dir)
			err = server.UploadToR2(srv.R2Client, baseDir, files)
		}()

		c.JSON(http.StatusOK, gin.H{
			"url":       req.RepoUrl,
			"projectId": projectId,
			"status":    "Clone in Progress",
		})
	}
}
