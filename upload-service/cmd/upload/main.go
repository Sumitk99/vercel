package main

import (
	"fmt"
	"github.com/Sumitk99/vercel/upload-service/routes"
	"github.com/Sumitk99/vercel/upload-service/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	PORT            string
	AccessKeyID     string
	SecretAccessKey string
	EndPoint        string
	RedisAddress    string
}

func main() {

	var router *gin.Engine = gin.New()
	var cfg Config
	err := godotenv.Load(".env")
	cfg.PORT = os.Getenv("PORT")
	cfg.EndPoint = os.Getenv("R2_ENDPOINT")
	cfg.AccessKeyID = os.Getenv("ACCESS_KEY_ID")
	cfg.SecretAccessKey = os.Getenv("SECRET_ACCESS_KEY")
	cfg.RedisAddress = os.Getenv("REDIS_ADDRESS")
	log.Println(cfg.RedisAddress)
	R2Client, err := server.ConnectToR2(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.EndPoint)
	if err != nil {
		log.Fatalf("Error connecting to R2: %s", err.Error())
	}
	RedisClient, err := server.ConnectToRedis(cfg.RedisAddress)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %s", err.Error())
	}
	router.Use(gin.Logger())

	corsPolicy := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	routes.SetupRoutes(router, &server.Server{
		R2Client:    R2Client,
		RedisClient: RedisClient,
	})

	router.Use(cors.New(corsPolicy))
	err = router.Run(fmt.Sprintf("0.0.0.0:%s", cfg.PORT))
	if err != nil {
		log.Println(err.Error())
	}
}
