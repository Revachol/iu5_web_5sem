package pkg

import (
	"fmt"

	"github.com/Revachol/iu5_web_5sem/internal/app/config"
	"github.com/Revachol/iu5_web_5sem/internal/app/handler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Application описывает основное приложение с конфигурацией, маршрутизатором и обработчиком запросов
type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

// NewApp создает новый экземпляр приложения
func NewApp(c *config.Config, r *gin.Engine, h *handler.Handler) *Application {
	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

// RunApp запускает веб-сервер и регистрирует маршруты и статические ресурсы
func (a *Application) RunApp() {
	logrus.Info("Server start up")

	// Регистрируем маршруты и статику
	a.Handler.RegisterHandler(a.Router)
	a.Handler.RegisterStatic(a.Router)

	// Формируем адрес сервера
	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)

	// Запуск сервера
	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Server down")
}
