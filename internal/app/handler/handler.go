package handler

import (
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) GetOrders(ctx *gin.Context) {
	var orders []repository.Order
	var err error

	searchQuery := ctx.Query("searchHistoricalObject")
	if searchQuery == "" {
		orders, err = h.Repository.GetOrders()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		orders, err = h.Repository.GetOrdersByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	estimate, err := h.Repository.GetEstimateData(1)
	var estimateCount int
	if err != nil {
		logrus.Error(err)
		estimateCount = 0
	} else {
		estimateCount = len(estimate.OrderIDs)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":                   time.Now().Format("15:04:05"),
		"historical_objects":     orders,
		"searchHistoricalObject": searchQuery,
		"estimate_count":         estimateCount,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "historical_object.html", gin.H{
		"historical_object": order,
	})
}

func (h *Handler) GetEstimate(ctx *gin.Context) {
	// получаем значение параметра id из query, а не из пути
	idStr := ctx.Query("id")
	if idStr == "" {
		logrus.Error("id parameter is missing")
		ctx.String(http.StatusBadRequest, "id parameter is missing")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "invalid id parameter")
		return
	}

	estimate, err := h.Repository.GetEstimateData(id)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, "could not get estimate data")
		return
	}

	var orders []repository.Order
	for _, orderID := range estimate.OrderIDs {
		order, err := h.Repository.GetOrder(orderID)
		if err != nil {
			logrus.Error(err)
			continue
		}
		orders = append(orders, order)
	}

	ctx.HTML(http.StatusOK, "historical_estimate.html", gin.H{
		"estimate_objects": orders,
		"estimate_id":      id,
	})
}
