package controllers

import (
	"github.com/Sumitk99/vercel/request-handler/constants"
	"github.com/Sumitk99/vercel/request-handler/server"
	"github.com/gin-gonic/gin"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type temp struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}

//func ReqController(srv *server.Server) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//host := c.Request.Host
//		//id := strings.Split(host, ".")[0]
//		//R2ObjectKey := filepath.Join(constants.OutputPath, id, c.Request.URL.Path)
//		var tem temp
//		err := c.BindJSON(&tem)
//		log.Println(tem)
//		R2ObjectKey := filepath.Join(constants.OutputPath, tem.ID, tem.URL)
//		log.Println(R2ObjectKey)
//		res, err := server.DownloadFileFromR2(srv.R2Client, R2ObjectKey)
//		if err != nil {
//			log.Println(err)
//			c.JSON(http.StatusNoContent, gin.H{
//				"error": err.Error(),
//			})
//		}
//		defer res.Body.Close()
//		contentType := mime.TypeByExtension(filepath.Ext(tem.URL))
//		if contentType == "" {
//			contentType = "application/octet-stream"
//		}
//
//		// Set response headers
//		c.Header("Content-Type", contentType)
//		c.Header("Cache-Control", "public, max-age=3600")
//
//		// Stream the file
//		c.DataFromReader(http.StatusOK, *res.ContentLength, contentType, res.Body, nil)
//
//	}
//}

func ReqController(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		}

		// Parse JSON payload
		if err := c.BindJSON(&reqData); err != nil {
			log.Println("Invalid request payload:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		log.Println("Received request:", reqData)

		// Ensure valid input
		if reqData.ID == "" || reqData.URL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID or URL"})
			return
		}

		// Handle Angular routing (if it's not a file request, serve index.html)
		requestedPath := reqData.URL
		if requestedPath == "/" || !strings.Contains(filepath.Ext(requestedPath), ".") {
			requestedPath = "index.html"
		}

		// Construct the R2 object key
		r2ObjectKey := filepath.Join(constants.OutputPath, reqData.ID, requestedPath)

		log.Println("Fetching from R2:", r2ObjectKey)

		// Fetch file from R2
		res, err := server.DownloadFileFromR2(srv.R2Client, r2ObjectKey)
		if err != nil {
			log.Println("File not found in R2:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		defer res.Body.Close()

		// Detect content type
		contentType := mime.TypeByExtension(filepath.Ext(requestedPath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// Set headers
		c.Header("Content-Type", contentType)
		c.Header("Cache-Control", "public, max-age=3600")

		// Stream the file
		c.DataFromReader(http.StatusOK, *res.ContentLength, contentType, res.Body, nil)
	}
}
