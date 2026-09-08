package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"labs5sem-electrical-load/internal/app/repository"
)

// ServiceView - данные, которые мы передаём HTML-шаблону
//
// LikeCount вычисляется здесь по длине LikedBy.
type ServiceView struct {
	ID          int
	Name        string
	PowerKW     float64
	LoadA       float64
	Description string
	Image       string
	Video       string
	LikeCount   int
}

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// преобразуем данные repository в данные для страницы
// Именно здесь считаются лайки
func makeServiceViews(services []repository.Service) []ServiceView {
	var result []ServiceView

	for _, service := range services {

		loadA := service.PowerKW * 1000 / 220

		result = append(result, ServiceView{
			ID:          service.ID,
			Name:        service.Name,
			PowerKW:     service.PowerKW,
			LoadA:       loadA,
			Description: service.Description,
			Image:       service.Image,
			Video:       service.Video,
			LikeCount:   len(service.LikedBy),
		})
	}

	return result
}

// Обычное открытие ленты.
// Начинаем с первого прибора.
func (h *Handler) GetFeed(ctx *gin.Context) {
	services, err := h.Repository.GetServices()
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, "Ошибка получения услуг")
		return
	}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"services": makeServiceViews(services),
	})
}

// Открытие ленты с конкретного прибора.
//
// Например:
// /feed/8
//
// → лента начинается с микроволновки.
func (h *Handler) GetFeedFromID(ctx *gin.Context) {
	idString := ctx.Param("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "Некорректный ID")
		return
	}

	services, err := h.Repository.GetServicesFromID(id)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "Услуга не найдена")
		return
	}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"services": makeServiceViews(services),
	})
}

// Страница добавления.
func (h *Handler) GetAddPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "add.html", gin.H{})
}

// В ЛР1 данные ещё не сохраняем.
//
// Форма отправляется, но коллекция не меняется.
func (h *Handler) SubmitAddPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "add.html", gin.H{})
}

// Каталог + числовая фильтрация по мощности.
func (h *Handler) GetCatalog(ctx *gin.Context) {
	maxPowerString := ctx.Query("max_power")

	var services []repository.Service
	var err error

	filterError := ""

	if maxPowerString == "" {
		// Поле пустое → показываем всё.
		services, err = h.Repository.GetServices()
	} else {
		maxPower, parseErr := strconv.ParseFloat(maxPowerString, 64)

		if parseErr != nil || maxPower < 0 {
			filterError = "Введите корректную мощность."

			services, err = h.Repository.GetServices()
		} else {
			services, err = h.Repository.GetServicesByMaxPower(maxPower)
		}
	}

	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, "Ошибка получения услуг")
		return
	}

	ctx.HTML(http.StatusOK, "catalog.html", gin.H{
		"services":    makeServiceViews(services),
		"max_power":   maxPowerString,
		"filterError": filterError,
	})
}
