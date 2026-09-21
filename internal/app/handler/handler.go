package handler

import (
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

	draft, err := h.Repository.GetDraftElectricityConsumer()

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Черновик не найден",
		)

		return
	}

	ctx.HTML(
		http.StatusOK,
		"add.html",
		gin.H{
			"draft": draft,
		},
	)
}

// 3. КАТАЛОГ

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

	var electricityConsumers []repository.ElectricityConsumer
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

	ctx.HTML(
		http.StatusOK,
		"catalog.html",
		gin.H{
			"electricityConsumers": electricityConsumers,
			"min_power":            minPowerStr,
			"max_power":            maxPowerStr,
			"filterError":          filterError,
		},
	)
}
