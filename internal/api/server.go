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
		logrus.Error("Ошибка инициализации репозитория")
		return
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	// Подключаем HTML-шаблоны.
	r.LoadHTMLGlob("templates/*")

	// CSS, JS, изображения и видео.
	r.Static("/static", "./resources")

	// localhost:8080 → лента.
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/feed")
	})

	// Лента с первого элемента.
	r.GET("/feed", h.GetFeed)

	// Лента начиная с выбранной карточки.
	// Например /feed/8.
	r.GET("/feed/:id", h.GetFeedFromID)

	// Добавление.
	r.GET("/add", h.GetAddPage)
	r.POST("/add", h.SubmitAddPage)

	// Каталог.
	r.GET("/catalog", h.GetCatalog)

	err = r.Run(":8080")
	if err != nil {
		logrus.Error(err)
	}

	log.Println("Server down")
}
