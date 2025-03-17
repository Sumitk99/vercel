package main

import (
	"github.com/Sumitk99/vercel/deploy-service/controllers"
	"github.com/Sumitk99/vercel/deploy-service/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PORT            string
	AccessKeyID     string
	SecretAccessKey string
	EndPoint        string
	RedisAddress    string
}

func main() {

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

	controllers.StartRedisListener(&server.Server{
		R2Client:    R2Client,
		RedisClient: RedisClient,
	})
}
