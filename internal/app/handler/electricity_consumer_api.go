package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"electricity_consumers/internal/app/currentuser"
	"electricity_consumers/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET СПИСОК УСЛУГ С ФИЛЬТРАЦИЕЙ

func (h *Handler) GetElectricityConsumersAPI(
	ctx *gin.Context,
) {

	minPowerStr := strings.TrimSpace(
		ctx.DefaultQuery("min_power", "0"),
	)

	maxPowerStr := strings.TrimSpace(
		ctx.DefaultQuery("max_power", "5"),
	)

	minPower, err := strconv.ParseFloat(
		strings.ReplaceAll(minPowerStr, ",", "."),
		64,
	)

	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	maxPower, err := strconv.ParseFloat(
		strings.ReplaceAll(maxPowerStr, ",", "."),
		64,
	)

	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if minPower < 0 || maxPower < minPower {
		ctx.Status(http.StatusBadRequest)
		return
	}

	consumers, err :=
		h.Repository.GetElectricityConsumersByPowerRange(
			minPower,
			maxPower,
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	currentUserID :=
		currentuser.GetCurrentUserID()

	response :=
		ds.ToElectricityConsumerResponseList(
			consumers,
			currentUserID,
		)

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

// GET ЛЕНТА

func (h *Handler) GetFeedAPI(
	ctx *gin.Context,
) {

	idStr := ctx.Param("id")

	var consumerID int

	//если ID не указан, получаем первую опубликованную услугу
	if idStr == "" {

		firstID, err :=
			h.Repository.GetNextElectricityConsumerID(
				0,
			)

		if err != nil {
			logrus.Error(err)

			ctx.Status(
				http.StatusNotFound,
			)
			return
		}

		consumerID = firstID

	} else {

		id, err := strconv.Atoi(idStr)

		if err != nil || id <= 0 {
			ctx.Status(
				http.StatusBadRequest,
			)
			return
		}

		consumerID = id

		//если next=true, получаем следующую опубликованную услугу
		if ctx.Query("next") == "true" {

			nextID, err :=
				h.Repository.GetNextElectricityConsumerID(
					id,
				)

			if err != nil {
				logrus.Error(err)

				ctx.Status(
					http.StatusNotFound,
				)
				return
			}

			consumerID = nextID
		}
	}

	consumer, err :=
		h.Repository.GetElectricityConsumer(
			consumerID,
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusNotFound,
		)
		return
	}

	currentUserID :=
		currentuser.GetCurrentUserID()

	response :=
		ds.ToElectricityConsumerResponse(
			consumer,
			currentUserID,
		)

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

// GET ЧЕРНОВИК

func (h *Handler) GetDraftAPI(
	ctx *gin.Context,
) {

	currentUserID :=
		currentuser.GetCurrentUserID()

	draft, err :=
		h.Repository.GetDraftElectricityConsumer(
			currentUserID,
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	if draft == nil {
		ctx.Status(
			http.StatusNotFound,
		)
		return
	}

	response :=
		ds.ToElectricityConsumerResponse(
			*draft,
			currentUserID,
		)

	ctx.JSON(
		http.StatusOK,
		response,
	)
}


func getUploadedFileContentType(
	header *multipart.FileHeader,
) (string, error) {

	file, err := header.Open()

	if err != nil {
		return "",
			fmt.Errorf(
				"не удалось открыть файл: %w",
				err,
			)
	}

	defer file.Close()

	buffer := make([]byte, 512)

	n, err := file.Read(buffer)

	if err != nil && err != io.EOF {
		return "",
			fmt.Errorf(
				"не удалось прочитать файл: %w",
				err,
			)
	}

	contentType :=
		http.DetectContentType(
			buffer[:n],
		)

	return contentType, nil
}

// POST СОЗДАНИЕ УСЛУГИ

func (h *Handler) CreateElectricityConsumerAPI(
	ctx *gin.Context,
) {

	err := ctx.Request.ParseMultipartForm(
		32 << 20,
	)

	if err != nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	name :=
		strings.TrimSpace(
			ctx.PostForm("name"),
		)

	if name == "" {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	currentUserID :=
		currentuser.GetCurrentUserID()

	//у пользователя может быть не более одного черновика
	oldDraft, err :=
		h.Repository.GetDraftElectricityConsumer(
			currentUserID,
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	if oldDraft != nil {
		ctx.Status(
			http.StatusConflict,
		)
		return
	}

	// ИЗОБРАЖЕНИЕ

	imageHeader, err :=
		ctx.FormFile("image")

	if err != nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	imageContentType, err :=
		getUploadedFileContentType(
			imageHeader,
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if !strings.HasPrefix(
		imageContentType,
		"image/",
	) {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	// ВИДЕО

	videoHeader, err :=
		ctx.FormFile("video")

	if err != nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	videoContentType, err :=
		getUploadedFileContentType(
			videoHeader,
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if !strings.HasPrefix(
		videoContentType,
		"video/",
	) {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	// ЗАГРУЗКА ИЗОБРАЖЕНИЯ В MINIO

	imageName, err :=
		h.Repository.UploadMedia(
			imageHeader,
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	// ЗАГРУЗКА ВИДЕО В MINIO

	videoName, err :=
		h.Repository.UploadMedia(
			videoHeader,
		)

	if err != nil {
		logrus.Error(err)

		//картинка уже была загружена, поэтому удаляем её
		_ = h.Repository.DeleteMedia(
			imageName,
		)

		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	// СОЗДАНИЕ ЗАПИСИ В POSTGRESQL

	draft, err :=
		h.Repository.CreateDraftWithMedia(
			currentUserID,
			name,
			imageName,
			videoName,
		)

	if err != nil {
		logrus.Error(err)

		_ = h.Repository.DeleteMedia(
			imageName,
		)

		_ = h.Repository.DeleteMedia(
			videoName,
		)

		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	response :=
		ds.ToElectricityConsumerResponse(
			draft,
			currentUserID,
		)

	ctx.JSON(
		http.StatusCreated,
		response,
	)
}

// PUT ПУБЛИКАЦИЯ УСЛУГИ

func (h *Handler) PublishElectricityConsumerAPI(
	ctx *gin.Context,
) {

	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(
		idStr,
		10,
		64,
	)

	if err != nil || id == 0 {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	var request ds.PublishElectricityConsumerRequest

	if err := ctx.ShouldBindJSON(
		&request,
	); err != nil {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	request.Description =
		strings.TrimSpace(
			request.Description,
		)

	if request.Description == "" {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if request.PowerKW <= 0 {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	currentUserID :=
		currentuser.GetCurrentUserID()

	err = h.Repository.PublishDraft(
		uint(id),
		currentUserID,
		request.Description,
		request.PowerKW,
	)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	consumer, err :=
		h.Repository.GetElectricityConsumer(
			int(id),
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	response :=
		ds.ToElectricityConsumerResponse(
			consumer,
			currentUserID,
		)

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

// DELETE УСЛУГИ

func (h *Handler) DeleteElectricityConsumerAPI(
	ctx *gin.Context,
) {

	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(
		idStr,
		10,
		64,
	)

	if err != nil || id == 0 {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	currentUserID :=
		currentuser.GetCurrentUserID()

	err = h.Repository.SoftDeleteOwnElectricityConsumer(
		uint(id),
		currentUserID,
	)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusNotFound,
		)
		return
	}

	ctx.Status(
		http.StatusOK,
	)
}

// POST LIKE

func (h *Handler) SetLikeAPI(
	ctx *gin.Context,
) {

	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(
		idStr,
		10,
		64,
	)

	if err != nil || id == 0 {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	var request ds.LikeRequest

	if err := ctx.ShouldBindJSON(
		&request,
	); err != nil {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if request.Value == nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if *request.Value != 0 &&
		*request.Value != 1 {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	currentUserID :=
		currentuser.GetCurrentUserID()

	err = h.Repository.SetLike(
		currentUserID,
		uint(id),
		*request.Value,
	)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	consumer, err :=
		h.Repository.GetElectricityConsumer(
			int(id),
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	response :=
		ds.ToElectricityConsumerResponse(
			consumer,
			currentUserID,
		)

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
