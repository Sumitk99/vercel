package routes

import (
	"github.com/Sumitk99/vercel/upload-service/controllers"
	"github.com/Sumitk99/vercel/upload-service/server"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, srv *server.Server) {
	router.GET("/deploy", controllers.Controller(srv))
	router.GET("/status/:projectId", controllers.FetchStatus(srv))
}
