
package handler

import (
	"electricity_consumers/internal/app/repository"

	"github.com/redis/go-redis/v9"
)

// Связывает HTTP-запросы с repository и Redis
type Handler struct {
	Repository *repository.Repository
	Redis      *redis.Client
}

// Создаёт новый Handler
func NewHandler(
	r *repository.Repository,
	redisClient *redis.Client,
) *Handler {
	return &Handler{
		Repository: r,
		Redis:      redisClient,
	}
}
