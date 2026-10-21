package handler

import (
	"net/http"
	"strings"

	"electricity_consumers/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


 // RegisterUserAPI godoc
 // @Summary Регистрация пользователя
 // @Description Создание нового пользователя с логином и паролем
 // @Tags Пользователи
 // @Accept json
 // @Produce json
 // @Param user body ds.RegisterUserRequest true "Данные пользователя"
 // @Success 201 {object} ds.UserResponse
 // @Failure 400 {string} string "Некорректный запрос"
 // @Failure 409 {string} string "Логин уже существует"
 // @Router /users [post]
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
