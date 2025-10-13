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

// DELETE /api/delete/historical_estimate/:estimate_id/historical_objects/:object_id
func (h *Handler) DeleteHObjectFromHEstimateAPI(ctx *gin.Context) {
	estimateIDStr := ctx.Param("estimate_id")
	objectIDStr := ctx.Param("object_id")

	estimateID, err := strconv.Atoi(estimateIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid estimate id parameter: %w", err))
		return
	}

	objectID, err := strconv.Atoi(objectIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid object id parameter: %w", err))
		return
	}

	err = h.Repository.DeleteHObjectFromHEstimate(estimateID, objectID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("failed to delete object from estimate: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"estimate_id": estimateID,
		"object_id":   objectID,
		"message":     "Object removed from estimate",
	})
}

// PUT /api/estimate/historical_objects/:estimate_id/:object_id/quantity
func (h *Handler) UpdateQuantityAPI(ctx *gin.Context) {
	estimateIDStr := ctx.Param("estimate_id")
	objectIDStr := ctx.Param("object_id")

	estimateID, err := strconv.Atoi(estimateIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid estimate id parameter: %w", err))
		return
	}

	objectID, err := strconv.Atoi(objectIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid object id parameter: %w", err))
		return
	}

	var req struct {
		Quantity *float64 `json:"quantity"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректное тело запроса"))
		return
	}

	if err := h.Repository.UpdateQuantity(estimateID, objectID, *req.Quantity); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("failed to update quantity: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"estimate_id": estimateID,
		"object_id":   objectID,
		"quantity":    req.Quantity,
		"message":     "Quantity updated successfully",
	})
}
