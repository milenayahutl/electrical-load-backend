package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"electricity_consumers/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


// LoginAPI godoc
// @Summary Аутентификация пользователя
// @Description Проверка логина и пароля с выдачей JWT
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param credentials body ds.LoginRequest true "Логин и пароль"
// @Success 200 {object} ds.LoginResponse
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неверный логин или пароль"
// @Router /auth/login [post]
func (h *Handler) LoginAPI(ctx *gin.Context) {

	var request ds.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	request.Login = strings.TrimSpace(request.Login)

	if request.Login == "" || request.Password == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	// Поиск пользователя по логину
	user, err := h.Repository.GetUserByLogin(request.Login)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	)

	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// Секретный ключ JWT
	secret := os.Getenv("JWT_SECRET")

	if len(secret) < 32 {
		logrus.Error("JWT_SECRET is missing or too short")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Срок действия JWT
	now := time.Now()
	expiresIn := 2 * time.Hour

	// Формируем JWT
	claims := ds.JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: strconv.FormatUint(
				uint64(user.ID), 10,
			),
			Issuer:    "powercheck",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			ID:        uuid.NewString(),
		},
	}

	// Подписываем JWT
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		[]byte(secret),
	)

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Возвращаем JWT пользователю
	ctx.JSON(
		http.StatusOK,
		ds.LoginResponse{
			AccessToken: tokenString,
			TokenType:   "Bearer",
			ExpiresIn:   int(expiresIn.Seconds()),
		},
	)
}


 // LogoutAPI godoc
 // @Summary Выход из системы
 // @Description Отзыв текущего JWT через Redis blacklist. Ответ без тела.
 // @Tags Аутентификация
 // @Security BearerAuth
 // @Success 200 {string} string "Успешный выход"
 // @Failure 401 {string} string "Требуется авторизация"
 // @Failure 503 {string} string "Redis недоступен"
 // @Router /auth/logout [post]
func (h *Handler) LogoutAPI(ctx *gin.Context) {

	// Получаем JWT, который проверил middleware
	value, exists := ctx.Get(jwtClaimsContextKey)

	if !exists {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	claims, ok := value.(*ds.JWTClaims)

	if !ok || claims == nil || claims.ExpiresAt == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	if h.Redis == nil {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	// Оставшееся время действия токена
	ttl := time.Until(claims.ExpiresAt.Time)

	if ttl <= 0 {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// Сохраняем ID JWT в Redis blacklist
	key := "jwt:blacklist:" + claims.ID

	// Значение содержит ID пользователя
	userID := strconv.FormatUint(
		uint64(claims.UserID), 10,
	)

	err := h.Redis.Set(
		ctx.Request.Context(),
		key,
		userID,
		ttl,
	).Err()

	if err != nil {
		logrus.Error("Redis logout error: ", err)
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	ctx.Status(http.StatusOK)
}
