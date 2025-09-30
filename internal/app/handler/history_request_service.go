package handler

import (
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

	ctx.HTML(http.StatusOK, "historical_estimate.html", gin.H{
		"estimate_objects": services,
		"estimate_id":      requestID,
	})
}
