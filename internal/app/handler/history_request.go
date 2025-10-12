package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetHistoricalRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	logrus.Info("Received ID:", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	historical_object, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.HTML(http.StatusOK, "historical_object.html", gin.H{
		"historical_object": historical_object,
	})
}

// GET /api/historical_estimate_count
func (h *Handler) GetDraftRequestAPI(ctx *gin.Context) {
	userID := 1
	order, err := h.Repository.GetDraftRequest(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if order.ID == 0 {
		// Если черновика нет — возвращаем пустую корзину
		ctx.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"orderID":   0,
			"itemCount": 0,
		})
		return
	}

	// Считаем количество услуг
	count := h.Repository.GetCartCount()

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"orderID":   order.ID,
		"itemCount": count,
	})
}

// GET /api/historical_estimate
func (h *Handler) GetAllHistoricalEstimateAPI(ctx *gin.Context) {
	status := ctx.Query("status")
	start := ctx.Query("start") // формат YYYY-MM-DD
	end := ctx.Query("end")

	orders, err := h.Repository.GetOrdersFiltered(status, start, end)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"orders": orders,
	})
}

// GET /api/historical_estimate/:id
func (h *Handler) GetHistoricalEstimateAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	// Получаем заказ и обхекты в нем
	order, materials, err := h.Repository.GetHistoricalRequest(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Конвертация sql.NullTime в *time.Time для DateForm и DateFinish
	var dateForm *time.Time
	if !order.SubmittedAt.IsZero() {
		dateForm = order.SubmittedAt
	}

	var dateFinish *time.Time
	if order.CompletedAt != nil {
		dateFinish = order.CompletedAt
	}

	// Формируем ответ
	orderResp := ds.OrderResponse{
		ID:         order.ID,
		Status:     order.Status,
		DateCreate: order.CreatedAt,
		DateForm:   dateForm,
		DateFinish: dateFinish,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"order":     orderResp,
		"materials": materials,
	})
}
