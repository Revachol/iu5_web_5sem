package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetHistoricalEstimate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("invalid request id parameter: ", err)
		ctx.String(http.StatusBadRequest, "invalid request id parameter")
		return
	}

	services, err := h.Repository.GetServicesByRequestID(requestID)
	if err != nil {
		logrus.Error("could not get services for request: ", err)
		ctx.String(http.StatusInternalServerError, "could not get services for request")
		return
	}

	var totalSum float64
	for _, service := range services {
		totalSum += service.Service.PriceUSD
	}

	ctx.HTML(http.StatusOK, "historical_estimate.html", gin.H{
		"estimate_objects": services,
		"estimate_id":      requestID,
		"total_sum":        totalSum,
		"total_cost_usd":   fmt.Sprintf("%.2f", totalSum*73.5), // форматирование до двух знаков после запятой
	})
}

func (h *Handler) DeleteRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("invalid request id parameter: ", err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	err = h.Repository.UpdateRequestStatus(requestID, "deleted")
	if err != nil {
		logrus.Error("failed to update request status: ", err)
	}

	ctx.Redirect(http.StatusFound, "/")
}
