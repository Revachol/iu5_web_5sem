package repository

import (
	"fmt"
	"math"
	"time"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
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

func (r *Repository) GetOrdersFiltered(status, start, end string) ([]ds.OrderResponse, error) {
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
	// Проверяем, что все wall_length заполнены
	var count int64
	if err := r.db.Model(&ds.Historical_request_service{}).
		Where("request_id = ? AND quantity IS NULL", estimateID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("ошибка проверки quantity: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("нельзя сформировать заказ: не все quantity заполнены")
	}

	// Обновляем заказ: статус и date_form
	updates := map[string]interface{}{
		"status":       "submitted",
		"submitted_at": time.Now(),
	}

	// Обновляем только если текущий статус черновик
	if err := r.db.Model(&ds.Historical_request{}).
		Where("id = ? AND status = ?", estimateID, "draft").
		Updates(updates).Error; err != nil {
		return fmt.Errorf("ошибка обновления заказа: %w", err)
	}

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
	const periodCastSQL = "CAST(HistoricalService.historical_period AS DECIMAL)"
	err := r.db.Model(&ds.Historical_request_service{}).
		Select("AVG("+periodCastSQL+")").
		Joins("JOIN historical_services AS HistoricalService ON historical_request_services.service_id = HistoricalService.id").
		Where("historical_request_services.request_id = ?", id).
		Scan(&avgYear).Error
	if err != nil {
		return 0, err
	}
	return int(avgYear), nil
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
	var request ds.Historical_request
	if err := r.db.First(&request, estimateID).Error; err != nil {
		return fmt.Errorf("заказ с ID=%d не найден", estimateID)
	}
	averageYear, _ := r.GetAverageServiceYear(estimateID)
	inflationMultiplier := r.СalculateInflationMultiplier(averageYear, 2025)
	totalHistoricalCost := r.CalculateTotalHistoricalCost(estimateID)
	finalCost := totalHistoricalCost * inflationMultiplier
	updates := map[string]interface{}{
		"status":               req.Status,
		"moderator_id":         req.ModeratorID,
		"completed_at":         time.Now(),
		"total_cost_usd":       finalCost,
		"inflation_multiplier": inflationMultiplier,
	}
	if err := r.db.Model(&ds.Historical_request{}).Where("id = ?", estimateID).Updates(updates).Error; err != nil {
		return fmt.Errorf("не удалось обновить итоговую стоимость: %w", err)
	}

	return nil

}
