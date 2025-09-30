package handler

import (
	"net/http"
	"strconv"

	"time"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	estimateCount := h.Repository.GetCartCount()
	draftRequestID, _ := h.Repository.GetDraftRequestID() // если нет — будет 0

	logrus.Info("Draft Request ID:", draftRequestID)
	logrus.Info("Estimate Count:", estimateCount)
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":                   time.Now().Format("15:04:05"),
		"historical_objects":     horders,
		"searchHistoricalObject": searchQuery,
		"estimate_count":         estimateCount,
		"draft_request_id":       draftRequestID,
	})
}

func (h *Handler) GetHistoricalObject(ctx *gin.Context) {
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

func (h *Handler) AddServiceToRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("invalid service id: ", err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	err = h.Repository.AddServiceToRequest(id)
	if err != nil {
		logrus.Error("failed to add service to request: ", err)
	}

	ctx.Redirect(http.StatusFound, "/")
}
