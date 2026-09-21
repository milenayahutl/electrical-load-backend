package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"electricity_consumers/internal/app/handler"
	"electricity_consumers/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error(err)
		return
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/electricity_consumers/feed/1")
	})

	r.GET("/electricity_consumers/feed", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/electricity_consumers/feed/1")
	})

	r.GET("/electricity_consumers/feed/:id", h.GetFeed)

	r.GET("/electricity_consumers/add", h.GetAddPage)

	r.GET("/electricity_consumers/catalog", h.GetCatalog)

	if err := r.Run(":8080"); err != nil {
		logrus.Error(err)
		return
	}

	log.Println("Server down")
}
