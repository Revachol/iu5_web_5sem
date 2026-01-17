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

	// Получаем user_id из контекста (если пользователь авторизован)
	userID, exists := ctx.Get("user_id")
	var estimateCount int64
	var draftRequestID int

	if exists {
		estimateCount, _ = h.Repository.GetCartCountForUser(userID.(int))
		draftRequestID, _ = h.Repository.GetDraftRequestID(userID.(int))
	} else {
		estimateCount = 0
		draftRequestID = 0
	}

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

// GET /api/historical_object/:id

// @Summary Получить услугу по ID
// @Description Возвращает полную информацию об исторической услуге по её ID.
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "ID исторической услуги (объекта)"
// @Success 200 {object} object "Успешный ответ"
// @Failure 400 {object} object "Неверный ID параметра"
// @Failure 404 {object} object "Объект не найден"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/historical_object/{id} [get]
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

// @Summary Получить список услуг
// @Description Возвращает список всех исторических услуг. Может быть отфильтрован по названию.
// @Tags services
// @Accept json
// @Produce json
// @Param title query string false "Название услуги (частичное совпадение)"
// @Success 200 {object} object "Список исторических услуг"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/historical_objects [get]
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

// @Summary Создать новую историческую услугу
// @Description Создает новую запись исторической услуги в базе данных.
// @Tags services
// @Accept json
// @Produce json
// @Param Historical_service body ds.Historical_service true "Данные для создания услуги"
// @Success 200 {object} object "Успешно созданная услуга"
// @Failure 400 {object} object "Некорректный JSON или неверные данные"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/create_object [post]
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

// @Summary Обновить историческую услугу
// @Description Обновляет существующую услугу по её ID.
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "ID исторической услуги для обновления"
// @Param Historical_service body ds.Historical_service true "Данные для обновления услуги"
// @Success 200 {object} object "Успешно обновленная услуга"
// @Failure 400 {object} object "Неверный ID или некорректный JSON"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/historical_object/{id} [put]
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

// @Summary Удалить историческую услугу
// @Description Удаляет услугу из базы данных по её ID.
// @Tags services
// @Produce json
// @Param id path int true "ID исторической услуги для удаления"
// @Success 200 {object} object "Сообщение об успешном удалении"
// @Failure 400 {object} object "Неверный ID параметра"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /api/historical_object/{id} [delete]
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

// @Summary Добавить услугу в черновую заявку
// @Description Добавляет указанную услугу в текущую черновую заявку (корзину) пользователя. Если черновика нет, он создается.
// @Tags estimates
// @Accept json
// @Produce json
// @Param id path int true "ID исторической услуги для добавления"
// @Success 200 {object} object "Сообщение об успешном добавлении и информация о заявке"
// @Failure 400 {object} object "Неверный ID параметра"
// @Failure 500 {object} object "Ошибка сервера (например, при создании черновика)"
// @Router /api/add_to_estimate/{id} [post]
func (h *Handler) AddHistoricalObjecsToRequestAPI(ctx *gin.Context) {
	objectIDStr := ctx.Param("id")
	objectID, err := strconv.Atoi(objectIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user_id not found in context"))
		return
	}

	// Получаем черновой заказ пользователя
	order, err := h.Repository.GetDraftRequest(userID.(int))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Если чернового заказа нет — создаём новый
	if order.ID == 0 {
		order, err = h.Repository.CreateDraftRequest(userID.(int))
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	// Добавляем материал в заказ
	if err := h.Repository.AddServiceToRequest(objectID, userID.(int)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Получаем новое количество материалов в заказе
	count, err := h.Repository.GetCartCountForUser(userID.(int))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Объект добавлен в черновой заказ",
		"estimateID": order.ID,
		"itemCount":  count,
	})
}

// POST /api/historical_object/:id/image

// @Summary Загрузить изображение для услуги
// @Description Загружает изображение для конкретной исторической услуги. Использует multipart/form-data.
// @Tags services
// @Accept mpfd
// @Produce json
// @Param id path int true "ID исторической услуги"
// @Param image formData file true "Файл изображения"
// @Success 200 {object} object "Сообщение об успешной загрузке"
// @Failure 400 {object} object "Неверный ID или файл не предоставлен"
// @Failure 500 {object} object "Ошибка сервера при обработке файла"
// @Router /api/historical_object/{id}/image [post]
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
