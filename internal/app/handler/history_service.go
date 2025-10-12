package handler

import (
	"net/http"
	"strconv"

	"time"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetHistoricalObjects(ctx *gin.Context) {
	var horders []ds.Historical_service
	var err error

	searchQuery := ctx.Query("searchHistoricalObject")
	if searchQuery == "" {
		horders, err = h.Repository.GetHistoricalObjects()
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

// GET /api/historical_object/:id
func (h *Handler) GetHistoricalObjectAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if order.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":      "error",
			"description": "Объект не найден",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":            "success",
		"historical_object": order,
	})
}

// GET /api/historical_objects?title=<название>
func (h *Handler) GetHistoricalObjectsAPI(ctx *gin.Context) {
	title := ctx.Query("title")

	objects, err := h.Repository.GetOrdersByTitle(title)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":             "success",
		"historical_objects": objects,
	})
}

// POST /api/create_object
func (h *Handler) CreateHistoricalObjectAPI(ctx *gin.Context) {
	var input ds.Historical_service
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.Repository.CreateHistoricalObject(&input); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":            "success",
		"historical_object": input,
	})
}

// PUT /api/historical_object/:id
func (h *Handler) UpdateHistoricalObjectAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var input ds.Historical_service
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.UpdateHistoricalObject(id, &input)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":            "success",
		"historical_object": input,
	})
}

// DELETE /api/historical_object/:id
func (h *Handler) DeleteHistoricalObjectAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteHistoricalObject(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Объект успешно удалён",
	})
}

// POST /api/add_to_estimate/:id
func (h *Handler) AddHistoricalObjecsToRequestAPI(ctx *gin.Context) {
	objectIDStr := ctx.Param("id")
	objectID, err := strconv.Atoi(objectIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	userID := 1

	// Получаем черновой заказ пользователя
	order, err := h.Repository.GetDraftRequest(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Если чернового заказа нет — создаём новый
	if order.ID == 0 {
		order, err = h.Repository.GetDraftRequest(userID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Добавляем материал в заказ
	if err := h.Repository.AddServiceToRequest(objectID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Получаем новое количество материалов в заказе
	count := h.Repository.GetCartCount()

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Объект добавлен в черновой заказ",
		"orderID":   order.ID,
		"itemCount": count,
	})
}

// POST /api/historical_object/:id/image
func (h *Handler) UploadHistoricalObjectImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
		return
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no image file"})
		return
	}

	if err := h.Repository.UploadHistoricalObjectImage(id, fileHeader); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "image uploaded"})
}
