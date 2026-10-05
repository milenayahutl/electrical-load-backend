package handler

import (
	"net/http"
	"strings"

	"electricity_consumers/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// POST АУТЕНТИФИКАЦИЯ

func (h *Handler) LoginAPI(
	ctx *gin.Context,
) {

	var request ds.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	request.Login =
		strings.TrimSpace(request.Login)

	request.Password =
		strings.TrimSpace(request.Password)

	if request.Login == "" ||
		request.Password == "" {

		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
}

// POST ДЕАВТОРИЗАЦИЯ

func (h *Handler) LogoutAPI(
	ctx *gin.Context,
) {

	ctx.Status(http.StatusOK)
}
