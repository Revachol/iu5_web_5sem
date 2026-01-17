package repository

import (
	"fmt"
	"math"
	"time"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetDraftRequest(creatorID int) (ds.Historical_request, error) {
	var request ds.Historical_request
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").First(&request).Error
	if err == gorm.ErrRecordNotFound {
		// Создаем новый черновик, если не найден
		request = ds.Historical_request{
			CreatorID: creatorID,
			Status:    "draft",
		}
		err = r.db.Create(&request).Error
	}
	return request, err
}

func (r *Repository) GetHistoricalRequest(id int) (ds.Historical_request, []ds.Historical_request_service, error) {
	var request ds.Historical_request
	if err := r.db.Preload("Creator").Preload("Moderator").First(&request, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ds.Historical_request{}, nil, fmt.Errorf("заявка с ID=%d не найдена", id)
		}
		return ds.Historical_request{}, nil, err
	}

	var services []ds.Historical_request_service
	if err := r.db.Preload("Service").Where("request_id = ?", id).Find(&services).Error; err != nil {
		return ds.Historical_request{}, nil, err
	}

	return request, services, nil
}

// IsEstimateOwnedByUser проверяет, принадлежит ли заявка пользователю
func (r *Repository) IsEstimateOwnedByUser(estimateID, userID int) bool {
	var count int64
	err := r.db.Model(&ds.Historical_request{}).
		Where("id = ? AND creator_id = ?", estimateID, userID).
		Count(&count).Error

	return err == nil && count > 0
}

func (r *Repository) GetOrdersFiltered(status, start, end string, userID int) ([]ds.OrderResponse, error) {
	var orders []ds.OrderResponse

	query := r.db.
		Table("historical_requests ho").
		Select(`ho.id, 
		        ho.status, 
		        ho.created_at, 
		        ho.submitted_at, 
		        ho.completed_at, 
		        ho.moderator_id, 
		        ho.creator_id`).
		Joins("LEFT JOIN users u1 ON u1.id = ho.moderator_id").
		Joins("LEFT JOIN users u2 ON u2.id = ho.creator_id").
		Where("ho.status NOT IN ?", []string{"deleted", "draft"})

	// Если указан userID (не 0), фильтруем по пользователю
	if userID > 0 {
		query = query.Where("ho.creator_id = ?", userID)
	}

	// фильтр по статусу
	if status != "" {
		query = query.Where("ho.status = ?", status)
	}

	// фильтр по диапазону дат
	if start != "" && end != "" {
		query = query.Where("ho.created_at BETWEEN ? AND ?", start, end)
	}

	if err := query.Scan(&orders).Error; err != nil {
		return nil, fmt.Errorf("ошибка при получении заказов: %w", err)
	}

	return orders, nil
}

func (r *Repository) UpdateHistoricalRequest(estimateID int, req ds.UpdateOrderRequest) error {
	updates := make(map[string]interface{})

	if req.CurrentYear != nil {
		updates["current_year"] = *req.CurrentYear
	}

	if len(updates) == 0 {
		return nil // ничего менять не нужно
	}

	return r.db.Model(&ds.Historical_request{}).Where("id = ?", estimateID).Updates(updates).Error
}

func (r *Repository) FormMaterialOrder(estimateID int) error {
	logrus.Infof("=== Repository.FormMaterialOrder ID=%d ===", estimateID)

	var count int64
	if err := r.db.Model(&ds.Historical_request_service{}).
		Where("request_id = ? AND quantity IS NULL", estimateID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("ошибка проверки quantity: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("нельзя сформировать заказ: не все quantity заполнены")
	}

	// Проверяем текущий год в заявке
	var request ds.Historical_request
	if err := r.db.First(&request, estimateID).Error; err != nil {
		return fmt.Errorf("заявка не найдена: %w", err)
	}
	logrus.Infof("Current year in request: %d", request.CurrentYear)

	// Обновляем заказ: статус, date_form и current_year (если не указан)
	updates := map[string]interface{}{
		"status":       "submitted",
		"submitted_at": time.Now(),
	}

	// Если год не указан, ставим 2025
	if request.CurrentYear == 0 {
		updates["current_year"] = 2025
		logrus.Info("CurrentYear was 0, setting to 2025")
	}

	// Обновляем только если текущий статус черновик
	if err := r.db.Model(&ds.Historical_request{}).
		Where("id = ? AND status = ?", estimateID, "draft").
		Updates(updates).Error; err != nil {
		return fmt.Errorf("ошибка обновления заказа: %w", err)
	}

	logrus.Info("=== Repository.FormMaterialOrder SUCCESS ===")
	return nil
}

func (r *Repository) DeleteHistoricalEstimete(id int) error {
	updates := map[string]interface{}{
		"status":       "deleted",
		"submitted_at": time.Now(), // дата завершения
	}
	return r.db.Model(&ds.Historical_request{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) GetAverageServiceYear(id int) (int, error) {
	var avgYear float64
	// Игнорируем некорректные года: < 1800 (чтобы исключить отрицательные и очень старые)
	const periodCastSQL = "CAST(HistoricalService.historical_period AS DECIMAL)"
	err := r.db.Model(&ds.Historical_request_service{}).
		Select("AVG("+periodCastSQL+")").
		Joins("JOIN historical_services AS HistoricalService ON historical_request_services.service_id = HistoricalService.id").
		Where("historical_request_services.request_id = ? AND "+periodCastSQL+" >= 1800", id).
		Scan(&avgYear).Error
	if err != nil {
		return 0, err
	}

	result := int(avgYear)
	logrus.Infof("GetAverageServiceYear result: %d (filtered >= 1800)", result)

	// Если результат все еще некорректный, возвращаем 0
	if result < 1800 || result > 2100 {
		logrus.Warnf("Average year %d is out of valid range, returning 0", result)
		return 0, nil
	}

	return result, nil
}

func (r *Repository) СalculateInflationMultiplier(startYear, endYear int) float64 {
	if startYear >= endYear {
		return 1.0
	}

	// Используем среднюю годовую инфляцию ~3%
	years := endYear - startYear
	inflationRate := 1.03 // 3% в год

	return math.Pow(inflationRate, float64(years))
}

func (r *Repository) CalculateTotalHistoricalCost(id int) float64 {
	var total float64
	err := r.db.Model(&ds.Historical_request_service{}).
		Select("SUM(total_price_usd)").
		Where("historical_request_services.request_id = ?", id).
		Scan(&total).Error
	if err != nil {
		return 0
	}
	return total
}

func (r *Repository) FinishHistoricalEstimate(estimateID int, req ds.CompleteOrderRequest) error {
	logrus.Infof("=== Repository.FinishHistoricalEstimate ID=%d ===", estimateID)

	var request ds.Historical_request
	if err := r.db.First(&request, estimateID).Error; err != nil {
		return fmt.Errorf("заказ с ID=%d не найден", estimateID)
	}
	logrus.Infof("Request found. CurrentYear from DB: %d", request.CurrentYear)

	// Если current_year не указан, используем 2025
	currentYear := request.CurrentYear
	if currentYear == 0 {
		currentYear = 2025
		logrus.Warn("CurrentYear was 0, using default: 2025")
	}

	// Получаем средний год услуг (игнорируя года < 1800)
	averageYear, err := r.GetAverageServiceYear(estimateID)
	if err != nil {
		logrus.Warnf("Failed to get average year: %v, using current year", err)
		averageYear = currentYear
	}
	logrus.Infof("Average service year: %d", averageYear)

	// Если средний год не получен или некорректный, используем текущий год
	if averageYear == 0 || averageYear < 1800 || averageYear > currentYear {
		logrus.Warnf("Invalid average year %d, using current year %d", averageYear, currentYear)
		averageYear = currentYear
	}

	// Дополнительная защита: если разница лет слишком большая (> 50 лет),
	// ограничиваем её чтобы избежать огромных чисел
	yearDifference := currentYear - averageYear
	if yearDifference > 50 {
		logrus.Warnf("Year difference too large (%d years), limiting to 50 years", yearDifference)
		averageYear = currentYear - 50
	}

	inflationMultiplier := r.СalculateInflationMultiplier(averageYear, currentYear)
	totalHistoricalCost := r.CalculateTotalHistoricalCost(estimateID)
	finalCost := totalHistoricalCost * inflationMultiplier

	logrus.Infof("Calculation: avgYear=%d, currYear=%d, yearDiff=%d, inflMult=%.4f, baseCost=%.2f, finalCost=%.2f",
		averageYear, currentYear, currentYear-averageYear, inflationMultiplier, totalHistoricalCost, finalCost)

	logrus.Infof("Calculation: avgYear=%d, currYear=%d, inflMult=%.4f, baseCost=%.2f, finalCost=%.2f",
		averageYear, currentYear, inflationMultiplier, totalHistoricalCost, finalCost)

	updates := map[string]interface{}{
		"status":               req.Status,
		"moderator_id":         req.ModeratorID,
		"completed_at":         time.Now(),
		"total_cost_usd":       finalCost,
		"inflation_multiplier": inflationMultiplier,
	}

	// Если current_year был 0, обновляем его
	if request.CurrentYear == 0 {
		updates["current_year"] = currentYear
		logrus.Info("Adding current_year to updates")
	}

	logrus.Info("Updating database...")
	if err := r.db.Model(&ds.Historical_request{}).Where("id = ?", estimateID).Updates(updates).Error; err != nil {
		logrus.Errorf("Database update failed: %v", err)
		return fmt.Errorf("не удалось обновить итоговую стоимость: %w", err)
	}

	logrus.Info("=== Repository.FinishHistoricalEstimate SUCCESS ===")
	return nil
}

// UpdateCalculatedPrice обновляет расчетную стоимость заявки
func (r *Repository) UpdateCalculatedPrice(requestID int, calculatedPrice float64) error {
	// Обновляем total_cost_usd в заявке
	if err := r.db.Model(&ds.Historical_request{}).
		Where("id = ?", requestID).
		Update("total_cost_usd", calculatedPrice).Error; err != nil {
		return fmt.Errorf("failed to update calculated price: %w", err)
	}
	return nil
}
