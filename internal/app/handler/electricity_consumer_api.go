package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"electricity_consumers/internal/app/currentuser"
	"electricity_consumers/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


// GetElectricityConsumersAPI godoc
// @Summary Получение списка электрических приборов
// @Description Получение опубликованных услуг с фильтрацией по мощности
// @Tags Услуги
// @Produce json
// @Param min_power query number false "Минимальная мощность, кВт"
// @Param max_power query number false "Максимальная мощность, кВт"
// @Success 200 {array} ds.ElectricityConsumerResponse
// @Failure 400 {string} string "Некорректный диапазон"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /electricity_consumers [get]
func (h *Handler) GetElectricityConsumersAPI(ctx *gin.Context) {
	minStr := strings.ReplaceAll(strings.TrimSpace(ctx.DefaultQuery("min_power", "0")), ",", ".")
	maxStr := strings.ReplaceAll(strings.TrimSpace(ctx.DefaultQuery("max_power", "5")), ",", ".")

	minPower, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректное значение min_power"})
		return
	}

	maxPower, err := strconv.ParseFloat(maxStr, 64)
	if err != nil || minPower < 0 || maxPower < minPower {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный диапазон мощности"})
		return
	}

	consumers, err := h.Repository.GetElectricityConsumersByPowerRange(minPower, maxPower)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	userID := currentuser.GetCurrentUserID(ctx)
	ctx.JSON(http.StatusOK, ds.ToElectricityConsumerResponseList(consumers, userID))
}


 // GetFeedAPI godoc
 // @Summary Получение ленты
 // @Description Получение опубликованной услуги по ID или следующей услуги при next=true
 // @Tags Лента
 // @Produce json
 // @Param id path int true "ID услуги"
 // @Param next query boolean false "Перейти к следующей услуге"
 // @Success 200 {object} ds.ElectricityConsumerFeedResponse
 // @Failure 400 {string} string "Некорректный ID"
 // @Failure 404 {string} string "Услуга не найдена"
 // @Router /electricity_consumers/feed/{id} [get]
func (h *Handler) GetFeedAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	consumerID := id

	if ctx.Query("next") == "true" {
		consumerID, err = h.Repository.GetNextElectricityConsumerID(id)
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	}

	consumer, err := h.Repository.GetElectricityConsumer(consumerID)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	userID := currentuser.GetCurrentUserID(ctx)
	ctx.JSON(http.StatusOK, ds.ToElectricityConsumerFeedResponse(consumer, userID))
}


 // GetDraftAPI godoc
 // @Summary Получение черновика
 // @Description Получение черновика текущего авторизованного пользователя
 // @Tags Услуги
 // @Produce json
 // @Security BearerAuth
 // @Success 200 {object} ds.ElectricityConsumerResponse
 // @Failure 401 {string} string "Требуется авторизация"
 // @Failure 404 {string} string "Черновик не найден"
 // @Router /electricity_consumers/draft [get]
func (h *Handler) GetDraftAPI(ctx *gin.Context) {
	userID := currentuser.GetCurrentUserID(ctx)

	draft, err := h.Repository.GetDraftElectricityConsumer(userID)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if draft == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, ds.ToElectricityConsumerResponse(*draft, userID))
}

// Определение MIME-типа загруженного файла
func getUploadedFileContentType(header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	return http.DetectContentType(buffer[:n]), nil
}

// Дополнительная проверка MP4
// Иногда Go определяет MP4 как application/octet-stream
func hasMP4Signature(header *multipart.FileHeader) bool {
	if !strings.EqualFold(filepath.Ext(header.Filename), ".mp4") {
		return false
	}

	file, err := header.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	headerBytes := make([]byte, 12)

	if _, err := io.ReadFull(file, headerBytes); err != nil {
		return false
	}

	return string(headerBytes[4:8]) == "ftyp"
}


 // CreateElectricityConsumerAPI godoc
 // @Summary Создание новой услуги
 // @Description Создание черновика с изображением и видео. Создатель определяется по JWT.
 // @Tags Услуги
 // @Accept multipart/form-data
 // @Produce json
 // @Security BearerAuth
 // @Param name formData string true "Название электроприбора"
 // @Param image formData file true "Изображение"
 // @Param video formData file true "Видео"
 // @Success 201 {object} ds.ElectricityConsumerResponse
 // @Failure 400 {string} string "Некорректные данные"
 // @Failure 401 {string} string "Требуется авторизация"
 // @Failure 409 {string} string "Черновик уже существует"
 // @Router /electricity_consumers [post]
func (h *Handler) CreateElectricityConsumerAPI(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(32 << 20)

	if err != nil {
		logrus.WithError(err).Warn("MULTIPART ERROR")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось прочитать form-data. Проверь Content-Type: должен быть multipart/form-data с boundary",
		})
		return
	}

	name := strings.TrimSpace(ctx.PostForm("name"))

	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Поле name не заполнено",
		})
		return
	}

	userID := currentuser.GetCurrentUserID(ctx)

	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Требуется авторизация",
		})
		return
	}

	oldDraft, err := h.Repository.GetDraftElectricityConsumer(userID)

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if oldDraft != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "У пользователя уже есть черновик",
		})
		return
	}

	// ИЗОБРАЖЕНИЕ
	imageHeader, err := ctx.FormFile("image")

	if err != nil {
		logrus.WithError(err).Warn("IMAGE FILE ERROR")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось получить файл image (Body -> form-data -> File)",
		})
		return
	}

	imageType, err := getUploadedFileContentType(imageHeader)

	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось прочитать изображение",
		})
		return
	}

	if !strings.HasPrefix(imageType, "image/") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":         "Файл image не распознан как изображение",
			"detected_mime": imageType,
		})
		return
	}

	// ВИДЕО
	videoHeader, err := ctx.FormFile("video")

	if err != nil {
		logrus.WithError(err).Warn("VIDEO FILE ERROR")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось получить файл video (Body -> form-data -> File)",
		})
		return
	}

	videoType, err := getUploadedFileContentType(videoHeader)

	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось прочитать видео",
		})
		return
	}

	if !strings.HasPrefix(videoType, "video/") && !hasMP4Signature(videoHeader) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":         "Файл video не распознан как видео или MP4",
			"detected_mime": videoType,
		})
		return
	}

	// ЗАГРУЗКА ИЗОБРАЖЕНИЯ В MINIO
	imageName, err := h.Repository.UploadMedia(imageHeader)

	if err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось загрузить изображение в MinIO",
		})
		return
	}

	// ЗАГРУЗКА ВИДЕО В MINIO
	videoName, err := h.Repository.UploadMedia(videoHeader)

	if err != nil {
		logrus.Error(err)

		// Удаляем изображение при ошибке загрузки видео
		_ = h.Repository.DeleteMedia(imageName)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось загрузить видео в MinIO",
		})
		return
	}

	// СОЗДАНИЕ ЧЕРНОВИКА В POSTGRESQL
	draft, err := h.Repository.CreateDraftWithMedia(
		userID,
		name,
		imageName,
		videoName,
	)

	if err != nil {
		logrus.Error(err)

		_ = h.Repository.DeleteMedia(imageName)
		_ = h.Repository.DeleteMedia(videoName)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось сохранить черновик в PostgreSQL",
		})
		return
	}

	// Успешное создание черновика
	ctx.JSON(
		http.StatusCreated,
		ds.ToElectricityConsumerResponse(draft, userID),
	)
}


 // PublishElectricityConsumerAPI godoc
 // @Summary Публикация услуги
 // @Description Обновление описания и мощности с расчётом силы тока. Публикация доступна создателю.
 // @Tags Услуги
 // @Accept json
 // @Produce json
 // @Security BearerAuth
 // @Param id path int true "ID черновика"
 // @Param data body ds.PublishElectricityConsumerRequest true "Описание и мощность"
 // @Success 200 {object} ds.ElectricityConsumerResponse
 // @Failure 400 {string} string "Некорректный запрос или публикация невозможна"
 // @Failure 401 {string} string "Требуется авторизация"
 // @Router /electricity_consumers/{id} [put]
func (h *Handler) PublishElectricityConsumerAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil || id == 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var request ds.PublishElectricityConsumerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	request.Description = strings.TrimSpace(request.Description)

	if request.Description == "" || request.PowerKW <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	userID := currentuser.GetCurrentUserID(ctx)

	err = h.Repository.PublishDraft(
		uint(id),
		userID,
		request.Description,
		request.PowerKW,
	)

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	consumer, err := h.Repository.GetElectricityConsumer(int(id))

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(
		http.StatusOK,
		ds.ToElectricityConsumerResponse(consumer, userID),
	)
}

 // DeleteElectricityConsumerAPI godoc
 // @Summary Удаление своей услуги
 // @Description Логическое удаление опубликованной услуги создателем
 // @Tags Услуги
 // @Security BearerAuth
 // @Param id path int true "ID услуги"
 // @Success 200 {string} string "Успешное удаление"
 // @Failure 401 {string} string "Требуется авторизация"
 // @Failure 404 {string} string "Услуга не найдена или не принадлежит пользователю"
 // @Router /electricity_consumers/{id} [delete]
func (h *Handler) DeleteElectricityConsumerAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil || id == 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	userID := currentuser.GetCurrentUserID(ctx)

	err = h.Repository.SoftDeleteOwnElectricityConsumer(
		uint(id),
		userID,
	)

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Status(http.StatusOK)
}

 // SetLikeAPI godoc
 // @Summary Установка или снятие лайка
 // @Description Значение value=1 ставит лайк, value=0 снимает лайк
 // @Tags Услуги
 // @Accept json
 // @Produce json
 // @Security BearerAuth
 // @Param id path int true "ID услуги"
 // @Param data body ds.LikeRequest true "Значение лайка"
 // @Success 200 {object} ds.ElectricityConsumerResponse
 // @Failure 400 {string} string "Некорректные данные"
 // @Failure 401 {string} string "Требуется авторизация"
 // @Router /electricity_consumers/{id}/likes [post]
func (h *Handler) SetLikeAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil || id == 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var request ds.LikeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil || request.Value == nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if *request.Value != 0 && *request.Value != 1 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	userID := currentuser.GetCurrentUserID(ctx)

	err = h.Repository.SetLike(
		userID,
		uint(id),
		*request.Value,
	)

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	consumer, err := h.Repository.GetElectricityConsumer(int(id))

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(
		http.StatusOK,
		ds.ToElectricityConsumerResponse(consumer, userID),
	)
}
