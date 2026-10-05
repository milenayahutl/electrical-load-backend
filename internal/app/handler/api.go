package handler

import "github.com/gin-gonic/gin"

// регистрирует маршруты веб-сервиса
func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")

	api.GET(
		"/electricity_consumers",
		h.GetElectricityConsumersAPI,
	)

	// лента без указанного id
	api.GET(
		"/electricity_consumers/feed",
		h.GetFeedAPI,
	)

	// лента по конкретному id
	api.GET(
		"/electricity_consumers/feed/:id",
		h.GetFeedAPI,
	)

	// черновик текущего пользователя
	api.GET(
		"/electricity_consumers/draft",
		h.GetDraftAPI,
	)

	// создание новой услуги черновика
	api.POST(
		"/electricity_consumers",
		h.CreateElectricityConsumerAPI,
	)

	// публикация черновика
	api.PUT(
		"/electricity_consumers/:id",
		h.PublishElectricityConsumerAPI,
	)

	// логическое удаление своей услуги
	api.DELETE(
		"/electricity_consumers/:id",
		h.DeleteElectricityConsumerAPI,
	)

	// лайк
	api.POST(
		"/electricity_consumers/:id/likes",
		h.SetLikeAPI,
	)

	// регистрация пользователя
	api.POST(
		"/users",
		h.RegisterUserAPI,
	)

	// аутентификация
	api.POST(
		"/auth/login",
		h.LoginAPI,
	)

	// деавторизация
	api.POST(
		"/auth/logout",
		h.LogoutAPI,
	)
}
