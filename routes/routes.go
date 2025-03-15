package routes

import (
	"github.com/Sumitk99/vercel/controllers"
	"github.com/Sumitk99/vercel/server"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, srv *server.Server) {
	router.GET("/deploy", controllers.Controller(srv))
}
