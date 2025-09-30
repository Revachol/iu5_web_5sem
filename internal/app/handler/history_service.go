package handler

import (
	"net/http"
	"strconv"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func (h *Handler) GetOrders(ctx *gin.Context) {
	var horders []ds.Historical_service
	var err error

	searchQuery := ctx.Query("searchHistoricalObject")
	if searchQuery == "" {
		horders, err = h.Repository.GetOrders()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		horders, err = h.Repository.GetOrdersByTitle(searchQuery)
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
		"historical_objects":     horders,
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
