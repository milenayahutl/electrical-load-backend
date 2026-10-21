
package pkg

import (
	"fmt"

	"electricity_consumers/internal/app/config"
	"electricity_consumers/internal/app/handler"

	_ "electricity_consumers/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(
	c *config.Config,
	r *gin.Engine,
	h *handler.Handler,
) *Application {

	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server start up")

	// Регистрация методов API
	a.Handler.RegisterAPI(a.Router)

	// Страница Swagger
	a.Router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	serverAddress := fmt.Sprintf(
		"%s:%d",
		a.Config.ServiceHost,
		a.Config.ServicePort,
	)

	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Server down")
}
