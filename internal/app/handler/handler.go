package handler

import (
	"electricity_consumers/internal/app/ds"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"electricity_consumers/internal/app/repository"
)

// Handler связывает HTTP-запросы с методами repository
type Handler struct {
	Repository *repository.Repository
}

const currentUserID uint = 1

// NewHandler создаёт новый handler и передаёт ему repository
func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// 1. ЛЕНТА

// GetFeed показывает один конкретный электроприбор
func (h *Handler) GetFeed(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Превращаем строку в число
	id, err := strconv.Atoi(idStr)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusBadRequest,
			"Некорректный ID",
		)

		return
	}

	electricityConsumer, err := h.Repository.GetElectricityConsumer(id)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusNotFound,
			"Услуга не найдена",
		)

		return
	}

	nextID, err := h.Repository.GetNextElectricityConsumerID(id)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения следующей услуги",
		)

		return
	}

	ctx.HTML(
		http.StatusOK,
		"feed.html",
		gin.H{
			"electricityConsumer": electricityConsumer,
			"nextID":              nextID,
		},
	)
}

// 2. СТРАНИЦА ДОБАВЛЕНИЯ

// GetAddPage показывает черновик
func (h *Handler) GetAddPage(ctx *gin.Context) {
	draft, err := h.Repository.GetDraftElectricityConsumer(currentUserID)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения черновика",
		)

		return
	}

	firstID, err := h.Repository.GetNextElectricityConsumerID(0)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения услуги",
		)

		return
	}

	ctx.HTML(
		http.StatusOK,
		"add.html",
		gin.H{
			"draft":    draft,
			"hasDraft": draft != nil,
			"firstID":  firstID,
		},
	)
}

func (h *Handler) GetCatalog(ctx *gin.Context) {

	minPowerStr := strings.TrimSpace(
		ctx.DefaultQuery(
			"min_power",
			"0",
		),
	)

	maxPowerStr := strings.TrimSpace(
		ctx.DefaultQuery(
			"max_power",
			"3",
		),
	)

	minPower, minErr := strconv.ParseFloat(
		strings.ReplaceAll(
			minPowerStr,
			",",
			".",
		),
		64,
	)

	maxPower, maxErr := strconv.ParseFloat(
		strings.ReplaceAll(
			maxPowerStr,
			",",
			".",
		),
		64,
	)

	var electricityConsumers []ds.ElectricityConsumer
	var err error

	filterError := ""

	if minErr != nil ||
		maxErr != nil ||
		minPower < 0 ||
		maxPower < minPower {

		filterError =
			"Введите корректный диапазон мощности."

		electricityConsumers, err =
			h.Repository.GetPublishedElectricityConsumers()

	} else {

		electricityConsumers, err =
			h.Repository.GetElectricityConsumersByPowerRange(
				minPower,
				maxPower,
			)
	}

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения электропотребителей",
		)

		return
	}

	firstID, err := h.Repository.GetNextElectricityConsumerID(0)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения услуги",
		)

		return
	}

	ctx.HTML(
		http.StatusOK,
		"catalog.html",
		gin.H{
			"electricityConsumers": electricityConsumers,
			"min_power":            minPowerStr,
			"max_power":            maxPowerStr,
			"filterError":          filterError,
			"firstID":              firstID,
		},
	)
}

func (h *Handler) CreateDraft(ctx *gin.Context) {
	name := strings.TrimSpace(ctx.PostForm("name"))

	if name == "" {
		ctx.String(
			http.StatusBadRequest,
			"Введите название",
		)
		return
	}

	_, err := h.Repository.CreateDraft(
		currentUserID,
		name,
	)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electricity_consumers/add",
	)
}

func (h *Handler) PublishDraft(ctx *gin.Context) {
	id, err := strconv.ParseUint(
		ctx.PostForm("id"),
		10,
		64,
	)

	if err != nil {
		ctx.String(
			http.StatusBadRequest,
			"Некорректный ID",
		)
		return
	}

	powerKW, err := strconv.ParseFloat(
		strings.ReplaceAll(
			ctx.PostForm("power_kw"),
			",",
			".",
		),
		64,
	)

	if err != nil {
		ctx.String(
			http.StatusBadRequest,
			"Некорректная мощность",
		)
		return
	}

	description := strings.TrimSpace(
		ctx.PostForm("description"),
	)

	if description == "" {
		ctx.String(
			http.StatusBadRequest,
			"Введите описание",
		)
		return
	}

	if powerKW < 0 {
		ctx.String(
			http.StatusBadRequest,
			"Мощность не может быть отрицательной",
		)
		return
	}

	err = h.Repository.PublishDraft(
		uint(id),
		currentUserID,
		description,
		powerKW,
	)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		fmt.Sprintf(
			"/electricity_consumers/feed/%d",
			id,
		),
	)
}

func (h *Handler) DeleteElectricityConsumer(ctx *gin.Context) {
	id, err := strconv.ParseUint(
		ctx.PostForm("id"),
		10,
		64,
	)

	if err != nil {
		ctx.String(
			http.StatusBadRequest,
			"Некорректный ID",
		)
		return
	}

	err = h.Repository.DeleteElectricityConsumer(
		uint(id),
	)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusFound,
		"/electricity_consumers/catalog",
	)
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET(
		"/electricity_consumers/feed/:id",
		h.GetFeed,
	)

	router.GET(
		"/electricity_consumers/add",
		h.GetAddPage,
	)

	router.GET(
		"/electricity_consumers/catalog",
		h.GetCatalog,
	)

	router.POST(
		"/electricity_consumers/add",
		h.CreateDraft,
	)

	router.POST(
		"/electricity_consumers/publish",
		h.PublishDraft,
	)

	router.POST(
		"/electricity_consumers/delete",
		h.DeleteElectricityConsumer,
	)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}
