package api

import (
	"log"

	"github.com/Revachol/iu5_web_5sem/internal/app/handler"
	"github.com/Revachol/iu5_web_5sem/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository("config/database_config.json")
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", handler.GetHistoricalObjects)
	r.GET("/historical_object/:id", handler.GetHistoricalObject)
	// r.GET("/historical_estimate", handler.GetEstimate)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
