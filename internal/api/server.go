package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"labs5sem-electrical-load/internal/app/handler"
	"labs5sem-electrical-load/internal/app/repository"
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
		ctx.Redirect(http.StatusFound, "/feed/1")
	})

	r.GET("/feed", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/feed/1")
	})

	r.GET("/feed/:id", h.GetFeed)

	r.GET("/add", h.GetAddPage)

	r.GET("/catalog", h.GetCatalog)

	if err := r.Run(":8080"); err != nil {
		logrus.Error(err)
		return
	}

	log.Println("Server down")
}
