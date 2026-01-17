package repository

import "github.com/Revachol/iu5_web_5sem/internal/app/ds"

func (r *Repository) GetServicesByRequestID(requestID int) ([]ds.Historical_request_service, error) {
	var services []ds.Historical_request_service
	err := r.db.Preload("Service").Where("request_id = ?", requestID).Find(&services).Error
	return services, err
}

func (r *Repository) UpdateRequestStatus(requestID int, status string) error {
	return r.db.Model(&ds.Historical_request{}).Where("id = ?", requestID).Update("status", status).Error
}

func (r *Repository) DeleteHObjectFromHEstimate(requestID, objectID int) error {
	return r.db.Where("request_id = ? AND service_id = ?", requestID, objectID).Delete(&ds.Historical_request_service{}).Error
}

func (r *Repository) UpdateQuantity(requestID, objectID int, quantity float64) error {
	// Сначала получаем текущую запись для обновления total_price_usd
	var record ds.Historical_request_service
	if err := r.db.Where("request_id = ? AND service_id = ?", requestID, objectID).First(&record).Error; err != nil {
		return err
	}

	// Обновляем quantity и пересчитываем total_price_usd
	record.Quantity = quantity
	record.TotalPriceUSD = record.UnitPriceUSD * quantity

	return r.db.Save(&record).Error
}
