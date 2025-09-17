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

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		orders, err = h.Repository.GetOrders()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		orders, err = h.Repository.GetOrdersByTitle(searchQuery) // в ином случае ищем заказ по заголовку
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
		"time":               time.Now().Format("15:04:05"),
		"historical_objects": orders,
		"query":              searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
		"estimate_count": estimateCount,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"historical_object": order,
	})
}

func (h *Handler) GetEstimate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		return
	}

	estimate, err := h.Repository.GetEstimateData(id)
	if err != nil {
		logrus.Error(err)
		// Maybe render an error page
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

	ctx.HTML(http.StatusOK, "estimate.html", gin.H{
		"estimate_objects": orders,
		"estimate_id":      id,
	})
}
