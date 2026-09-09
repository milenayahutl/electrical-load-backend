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
// с методами repository.
type Handler struct {
	Repository *repository.Repository
}

// NewHandler создаёт новый handler
// и передаёт ему repository.
func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// =====================================================
// 1. ЛЕНТА
// =====================================================

// GetFeed показывает один конкретный электроприбор.
//
// Например:
// /feed/1
// /feed/4
// /feed/15
func (h *Handler) GetFeed(ctx *gin.Context) {

	// Получаем ID из URL.
	//
	// Например:
	// /feed/15
	//
	// ctx.Param("id") вернёт строку "15".
	idStr := ctx.Param("id")

	// Превращаем строку "15"
	// в число 15.
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

	// Получаем ID следующей существующей услуги.
	//
	// ВАЖНО:
	// здесь НЕ используется id + 1.
	//
	// Например, если есть:
	// 1, 4, 15, 115
	//
	// после 4 repository вернёт 15.
	nextID, err := h.Repository.GetNextServiceID(id)

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения следующей услуги",
		)

		return
	}

	// Нужен для ссылки "больше".
	//
	// /feed/1
	// showFull = false
	//
	// /feed/1?full=1
	// showFull = true
	showFull := ctx.Query("full") == "1"

	// Передаём данные в feed.html.
	//
	// Внутри service уже лежат:
	// ID
	// Name
	// PowerKW
	// Description
	// Image
	// Video
	// LikedBy
	//
	// Поэтому отдельно передавать
	// каждую характеристику не нужно.
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

// =====================================================
// 2. СТРАНИЦА ДОБАВЛЕНИЯ
// =====================================================

// GetAddPage показывает существующий черновик.
//
// Пользователь сможет изменить текст
// прямо в input/textarea,
// но эти изменения НЕ отправляются на сервер
// и НЕ сохраняются в repository.
func (h *Handler) GetAddPage(ctx *gin.Context) {

	// Получаем заранее существующий черновик.
	//
	// Например, у тебя это может быть блендер.
	draft, err := h.Repository.GetDraftService()

	if err != nil {
		logrus.Error(err)

		ctx.String(
			http.StatusInternalServerError,
			"Черновик не найден",
		)

		return
	}

	// Передаём ВЕСЬ объект draft в add.html.
	//
	// В нём уже находятся:
	//
	// draft.Name
	// draft.PowerKW
	// draft.Description
	// draft.Image
	// draft.Video
	// draft.LikedBy
	//
	// Поэтому никаких отдельных:
	//
	// "name": ...
	// "power": ...
	// "image": ...
	//
	// здесь уже не нужно.
	ctx.HTML(
		http.StatusOK,
		"add.html",
		gin.H{
			"draft": draft,
		},
	)
}

// =====================================================
// 3. КАТАЛОГ
// =====================================================

// GetCatalog показывает все услуги
// либо фильтрует их по максимальной мощности.
//
// Например:
//
// /catalog
//
// покажет все.
//
// А:
//
// /catalog?max_power=2.0
//
// покажет приборы с мощностью <= 2.0 кВт.
func (h *Handler) GetCatalog(ctx *gin.Context) {
	maxPowerStr := strings.TrimSpace(
		ctx.Query("max_power"),
	)

	var services []repository.Service
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
