package handler

import (
	"github.com/Revachol/iu5_web_5sem/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetHistoricalObjects)
	router.GET("/historical_object/:id", h.GetHistoricalObject)
	router.GET("/historical_estimate/:id", h.GetHistoricalEstimate)
	router.GET("/api/historical_object/:id", h.GetHistoricalObjectAPI) //GET одна запись
	router.GET("/api/historical_objects", h.GetHistoricalObjectsAPI)   //GET список с фильтрацией

	router.POST("/add_to_cart/:id", h.AddServiceToRequest)
	router.POST("/delete_request/:id", h.DeleteRequest)
	router.POST("/api/create_object", h.CreateHistoricalObjectAPI)
	router.POST("/api/add_to_estimate/:id", h.AddHistoricalObjecsToRequestAPI)
	router.POST("/api/historical_object/:id/image", h.UploadHistoricalObjectImage)

	router.PUT("/api/historical_object/:id", h.UpdateHistoricalObjectAPI)
	router.DELETE("/api/historical_object/:id", h.DeleteHistoricalObjectAPI)

}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
