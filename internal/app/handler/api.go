
package handler

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")

	api.GET(
		"/electricity_consumers",
		h.OptionalAuth(),
		h.GetElectricityConsumersAPI,
	)

	api.GET(
		"/electricity_consumers/feed/:id",
		h.OptionalAuth(),
		h.GetFeedAPI,
	)

	api.POST(
		"/users",
		h.RegisterUserAPI,
	)

	api.POST(
		"/auth/login",
		h.LoginAPI,
	)


	protected := api.Group("")
	protected.Use(h.RequireAuth())

	protected.GET(
		"/electricity_consumers/draft",
		h.GetDraftAPI,
	)

	protected.POST(
		"/electricity_consumers",
		h.CreateElectricityConsumerAPI,
	)

	protected.PUT(
		"/electricity_consumers/:id",
		h.PublishElectricityConsumerAPI,
	)

	protected.DELETE(
		"/electricity_consumers/:id",
		h.DeleteElectricityConsumerAPI,
	)

	protected.POST(
		"/electricity_consumers/:id/likes",
		h.SetLikeAPI,
	)

	protected.POST(
		"/auth/logout",
		h.LogoutAPI,
	)
}
