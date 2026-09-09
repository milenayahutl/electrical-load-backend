package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"labs5sem-electrical-load/internal/app/repository"
)

// Handler связывает HTTP-запросы
// с методами repository
type Handler struct {
	Repository *repository.Repository
}

// NewHandler создаёт новый handler
// и передаёт ему repository
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

	// Получаем саму услугу по её ID.
	service, err := h.Repository.GetService(id)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusNotFound,
			"Услуга не найдена",
		)

		return
	}

	nextID, err := h.Repository.GetNextServiceID(id)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения следующей услуги",
		)

		return
	}

	showFull := ctx.Query("full") == "1"

	ctx.HTML(
		http.StatusOK,
		"feed.html",
		gin.H{
			"service":  service,
			"nextID":   nextID,
			"showFull": showFull,
		},
	)
}

// 2. СТРАНИЦА ДОБАВЛЕНИЯ

// GetAddPage показывает существующий черновик.
func (h *Handler) GetAddPage(ctx *gin.Context) {

	draft, err := h.Repository.GetDraftService()

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

// GetCatalog показывает все услуги
func (h *Handler) GetCatalog(ctx *gin.Context) {
	maxPowerStr := strings.TrimSpace(
		ctx.Query("max_power"),
	)

	var services []repository.Appliance
	var err error

	filterError := ""

	if maxPowerStr == "" {

		services, err =
			h.Repository.GetPublishedServices()

	} else {

		maxPower, parseErr := strconv.ParseFloat(
			strings.ReplaceAll(
				maxPowerStr,
				",",
				".",
			),
			64,
		)

		if parseErr != nil || maxPower < 0 {

			filterError = "Введите корректную мощность."

			services, err =
				h.Repository.GetPublishedServices()

		} else {

			services, err =
				h.Repository.GetServicesByMaxPower(maxPower)
		}
	}

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения услуг",
		)

		return
	}

	ctx.HTML(
		http.StatusOK,
		"catalog.html",
		gin.H{
			"services":    services,
			"max_power":   maxPowerStr,
			"filterError": filterError,
		},
	)
}
