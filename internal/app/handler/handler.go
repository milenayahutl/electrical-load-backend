package handler

import "electricity_consumers/internal/app/repository"

// связывает HTTP-запросы с repository
type Handler struct {
	Repository *repository.Repository
}

// создаёт новый Handler
func NewHandler(
	r *repository.Repository,
) *Handler {

	return &Handler{
		Repository: r,
	}
}
