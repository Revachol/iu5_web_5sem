package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// URL Django сервиса расчетов
	CalculationServiceURL = "http://localhost:8000/api/calculate/"
	// Токен авторизации (должен совпадать с токеном в Django .env)
	CalculationAuthToken = "12345678"
	// URL для callback от Django
	CallbackURL = "http://localhost:8080/api/update_calculated_price"
	// Коэффициент инфляции по умолчанию (5%)
	DefaultInflationRate = 0.05
)

// CalculationService структура услуги для отправки в Django
type CalculationService struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
	Year     int     `json:"year"`
}

// CalculationRequest структура запроса в Django сервис
type CalculationRequest struct {
	RequestID       string               `json:"request_id"`
	Services        []CalculationService `json:"services"`
	InflationRate   float64              `json:"inflation_rate"`
	CalculationYear int                  `json:"calculation_year"`
	CallbackURL     string               `json:"callback_url"`
}

// CalculationResult структура результата от Django сервиса
type CalculationResult struct {
	RequestID       string  `json:"request_id"`
	CalculatedPrice float64 `json:"calculated_price"`
	Status          string  `json:"status"`
	ErrorMessage    string  `json:"error_message,omitempty"`
}

// SendCalculationRequest отправляет заявку на расчет в Django сервис
func (h *Handler) SendCalculationRequest(requestID int, services []ds.Historical_request_service, currentYear int) error {
	// Формируем список услуг для расчета
	calcServices := make([]CalculationService, len(services))
	for i, service := range services {
		// Извлекаем год из периода или используем текущий год
		year := extractYearFromPeriod(service.Service.HistoricalPeriod)
		if year == 0 {
			year = currentYear
		}

		calcServices[i] = CalculationService{
			Price:    service.UnitPriceUSD,
			Quantity: service.Quantity,
			Year:     year,
		}
	}

	// Формируем запрос
	request := CalculationRequest{
		RequestID:       strconv.Itoa(requestID),
		Services:        calcServices,
		InflationRate:   DefaultInflationRate,
		CalculationYear: currentYear,
		CallbackURL:     CallbackURL,
	}

	// Сериализуем в JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal calculation request: %w", err)
	}

	logrus.Infof("Sending calculation request for order %d to Django service", requestID)

	// Отправляем POST запрос
	resp, err := http.Post(
		CalculationServiceURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send calculation request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("calculation service returned status %d", resp.StatusCode)
	}

	logrus.Infof("Calculation request for order %d accepted by Django service", requestID)
	return nil
}

// UpdateCalculatedPriceAPI обрабатывает callback от Django сервиса с результатом расчета
// POST /api/update_calculated_price
func (h *Handler) UpdateCalculatedPriceAPI(ctx *gin.Context) {
	// Проверка авторизации
	authHeader := ctx.GetHeader("Authorization")
	expectedAuth := fmt.Sprintf("Bearer %s", CalculationAuthToken)

	if authHeader != expectedAuth {
		logrus.Warn("Unauthorized calculation callback attempt")
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	// Парсим результат
	var result CalculationResult
	if err := ctx.BindJSON(&result); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("Received calculation result for request %s: %.2f (status: %s)",
		result.RequestID, result.CalculatedPrice, result.Status)

	// Конвертируем request_id в int
	requestID, err := strconv.Atoi(result.RequestID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid request_id: %w", err))
		return
	}

	// Обновляем заявку в БД
	if result.Status == "success" {
		if err := h.Repository.UpdateCalculatedPrice(requestID, result.CalculatedPrice); err != nil {
			logrus.Errorf("Failed to update calculated price for request %d: %v", requestID, err)
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		logrus.Infof("Successfully updated calculated price %.2f for request %d",
			result.CalculatedPrice, requestID)
	} else {
		logrus.Errorf("Calculation failed for request %d: %s",
			requestID, result.ErrorMessage)
	}

	// Отправляем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Result received",
	})
}

// extractYearFromPeriod извлекает год из строки периода
// Примеры: "18 век" -> 1750, "1800-1850" -> 1825, "2020" -> 2020
func extractYearFromPeriod(period string) int {
	// Простая реализация - можно улучшить
	// Попробуем найти 4-значное число в строке
	var year int
	_, err := fmt.Sscanf(period, "%d", &year)
	if err == nil && year >= 1000 && year <= 9999 {
		return year
	}

	// Если не нашли год, возвращаем текущий год
	return time.Now().Year()
}
