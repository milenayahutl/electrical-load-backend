
package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"electricity_consumers/internal/app/config"
	"electricity_consumers/internal/app/dsn"
	"electricity_consumers/internal/app/handler"
	"electricity_consumers/internal/app/repository"
	"electricity_consumers/internal/pkg"
)


 // @title PowerCheck API
 // @version 1.0
 // @description Веб-сервис для учёта электрических приборов, расчёта электрического тока и публикации услуг.
 // @host localhost:8080
 // @BasePath /api
 // @schemes http

 // @securityDefinitions.apikey BearerAuth
 // @in header
 // @name Authorization
 // @description Введите Bearer и JWT-токен через пробел


func main() {
	router := gin.Default()

	// Загрузка конфигурации
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf(
			"error loading config: %v",
			err,
		)
	}

	// Подключение к PostgreSQL и MinIO
	postgresString := dsn.FromEnv()

	rep, err := repository.New(
		&repository.RepositorySettings{
			PostgresDSN: postgresString,

			MinioEndpoint: os.Getenv(
				"MINIO_ENDPOINT",
			),

			MinioAccessKey: os.Getenv(
				"MINIO_ACCESS_KEY",
			),

			MinioSecretKey: os.Getenv(
				"MINIO_SECRET_KEY",
			),

			MinioBucketName: os.Getenv(
				"MINIO_BUCKET_NAME",
			),
		},
	)

	if err != nil {
		logrus.Fatalf(
			"error initializing repository: %v",
			err,
		)
	}

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	defer redisClient.Close()

	// Проверка соединения с Redis
	redisCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	err = redisClient.Ping(redisCtx).Err()
	cancel()

	if err != nil {
		logrus.Fatalf(
			"error connecting to Redis: %v",
			err,
		)
	}

	logrus.Info("Redis connected successfully")

	// Создание обработчика
	hand := handler.NewHandler(
		rep,
		redisClient,
	)

	// Создание приложения
	application := pkg.NewApp(
		conf,
		router,
		hand,
	)

	application.RunApp()
}
