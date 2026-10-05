package handler

import (
	"net/http"
	"strings"

	"electricity_consumers/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// POST РЕГИСТРАЦИЯ ПОЛЬЗОВАТЕЛЯ

func (h *Handler) RegisterUserAPI(
	ctx *gin.Context,
) {

	var request ds.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	request.Login =
		strings.TrimSpace(request.Login)

	request.Password =
		strings.TrimSpace(request.Password)

	if request.Login == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if request.Password == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	user, err :=
		h.Repository.CreateUser(
			request.Login,
			request.Password,
		)

	if err != nil {
		logrus.Error(err)

		ctx.Status(http.StatusConflict)
		return
	}

	response :=
		ds.ToUserResponse(user)

	ctx.JSON(
		http.StatusCreated,
		response,
	)
}
