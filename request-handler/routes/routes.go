package routes

import (
	"github.com/Sumitk99/vercel/request-handler/controllers"
	"github.com/Sumitk99/vercel/request-handler/server"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, srv *server.Server) {
	router.GET("/*filepath", controllers.ReqController(srv))
}
