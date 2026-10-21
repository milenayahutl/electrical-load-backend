package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"electricity_consumers/internal/app/currentuser"
	"electricity_consumers/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// Ключ для хранения данных JWT в контексте запроса
const jwtClaimsContextKey = "jwt_claims"

// Проверка JWT и Redis blacklist
func (h *Handler) getJWTClaims(
	ctx *gin.Context,
) (*ds.JWTClaims, int) {

	// Получаем JWT из заголовка Authorization
	header := ctx.GetHeader("Authorization")

	if !strings.HasPrefix(header, "Bearer ") {
		return nil, http.StatusUnauthorized
	}

	tokenString := strings.TrimSpace(
		strings.TrimPrefix(header, "Bearer "),
	)

	if tokenString == "" {
		return nil, http.StatusUnauthorized
	}

	// Получаем секретный ключ
	secret := os.Getenv("JWT_SECRET")

	if len(secret) < 32 {
		logrus.Error("JWT_SECRET is missing or too short")
		return nil, http.StatusInternalServerError
	}

	// Проверяем подпись и срок действия JWT
	claims := &ds.JWTClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("powercheck"),
		jwt.WithExpirationRequired(),
	)

	if err != nil || !token.Valid {
		return nil, http.StatusUnauthorized
	}

	// Проверяем идентификатор пользователя
	if claims.UserID == 0 ||
		claims.ID == "" ||
		claims.ExpiresAt == nil ||
		claims.Subject != strconv.FormatUint(
			uint64(claims.UserID), 10,
		) {

		return nil, http.StatusUnauthorized
	}

	// Проверяем, не находится ли JWT в blacklist
	if h.Redis == nil {
		return nil, http.StatusServiceUnavailable
	}

	key := "jwt:blacklist:" + claims.ID

	exists, err := h.Redis.Exists(
		ctx.Request.Context(),
		key,
	).Result()

	if err != nil {
		logrus.Error("Redis blacklist check: ", err)
		return nil, http.StatusServiceUnavailable
	}

	// Если JWT находится в Redis, доступ запрещён
	if exists > 0 {
		return nil, http.StatusUnauthorized
	}

	return claims, 0
}

// Обязательная авторизация
func (h *Handler) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		claims, status := h.getJWTClaims(ctx)

		if status != 0 {
			ctx.AbortWithStatus(status)
			return
		}

		// Передаём ID пользователя обработчикам
		ctx.Set(currentuser.UserIDKey, claims.UserID)
		ctx.Set(jwtClaimsContextKey, claims)

		ctx.Next()
	}
}

// Необязательная авторизация для публичных GET
func (h *Handler) OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		claims, status := h.getJWTClaims(ctx)

		if status == 0 {
			ctx.Set(currentuser.UserIDKey, claims.UserID)
			ctx.Set(jwtClaimsContextKey, claims)
		}

		ctx.Next()
	}
}
