package handler

import (
	"fmt"
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

// GetAllHistoricalEstimateAPI возвращает список исторических заявок с фильтрацией по статусу и диапазону дат.
//
// @Summary Получить список исторических заявок
// @Description Возвращает список заявок (не удалённых и не черновиков), с возможностью фильтрации по статусу и диапазону дат формирования.
// @Tags Estimates
// @Accept json
// @Produce json
// @Param status query string false "Фильтр по статусу заявки (например: draft, completed, rejected)"
// @Param start query string false "Дата начала диапазона (формат: YYYY-MM-DD)"
// @Param end query string false "Дата конца диапазона (формат: YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Список заявок успешно получен"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 500 {object} map[string]string "Ошибка на стороне сервера"
// @Router /api/historical_estimate [get]
// @Security ApiKeyAuth
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

// PUT /api/historical_estimate/:id
func (h *Handler) UpdateHistoricalEstimateAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var raw map[string]interface{}
	if err := ctx.BindJSON(&raw); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	allowed := map[string]bool{
		"current_year": true,
	}

	for k := range raw {
		if !allowed[k] {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("недопустимое поле: %s", k))
			return
		}
	}
	var req ds.UpdateOrderRequest
	if v, ok := raw["current_year"]; ok {
		if f, ok := v.(float64); ok {
			req.CurrentYear = &f
		}
	}

	if err := h.Repository.UpdateHistoricalRequest(id, req); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	order, _, err := h.Repository.GetHistoricalRequest(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"order":  order,
	})
}

// PUT /api/historical_estimate/:id/form
func (h *Handler) FormHistoricalEstimateAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	estimateID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.FormMaterialOrder(estimateID); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Возвращаем обновлённый заказ
	estimate, _, err := h.Repository.GetHistoricalRequest(estimateID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":              "success",
		"historical_estimate": estimate,
	})

}

func (h *Handler) DeleteHistoricalEstimeteAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteHistoricalEstimete(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Объект успешно удалён",
	})
}

// PUT /api/historical_estimate/:id/complete
func (h *Handler) FinishHistoricalEstimateAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req ds.CompleteOrderRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	req.Status = "completed" // или "rejected"
	req.ModeratorID = 1

	if err := h.Repository.FinishHistoricalEstimate(id, req); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	order, materials, err := h.Repository.GetHistoricalRequest(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"order":     order,
		"materials": materials,
	})
}
