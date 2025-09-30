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
	router.GET("/", h.GetOrders)
	router.GET("/order/:id", h.GetOrder)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./styles")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

// func (h *Handler) GetEstimate(ctx *gin.Context) {
// 	// получаем значение параметра id из query, а не из пути
// 	idStr := ctx.Query("id")
// 	if idStr == "" {
// 		logrus.Error("id parameter is missing")
// 		ctx.String(http.StatusBadRequest, "id parameter is missing")
// 		return
// 	}
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		logrus.Error(err)
// 		ctx.String(http.StatusBadRequest, "invalid id parameter")
// 		return
// 	}

// 	estimate, err := h.Repository.GetEstimateData(id)
// 	if err != nil {
// 		logrus.Error(err)
// 		ctx.String(http.StatusInternalServerError, "could not get estimate data")
// 		return
// 	}

// 	var orders []repository.Order
// 	for _, orderID := range estimate.OrderIDs {
// 		order, err := h.Repository.GetOrder(orderID)
// 		if err != nil {
// 			logrus.Error(err)
// 			continue
// 		}
// 		orders = append(orders, order)
// 	}

// 	ctx.HTML(http.StatusOK, "historical_estimate.html", gin.H{
// 		"estimate_objects": orders,
// 		"estimate_id":      id,
// 	})
// }
