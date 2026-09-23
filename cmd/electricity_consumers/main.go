package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"electricity_consumers/internal/app/config"
	"electricity_consumers/internal/app/dsn"
	"electricity_consumers/internal/app/handler"
	"electricity_consumers/internal/app/repository"
	"electricity_consumers/internal/pkg"
)

func main() {
	router := gin.Default()

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()

	rep, err := repository.New(postgresString)
	if err != nil {
		logrus.Fatalf(
			"error initializing repository: %v",
			err,
		)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(
		conf,
		router,
		hand,
	)

	application.RunApp()
}
