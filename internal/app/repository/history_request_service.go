package repository

import "github.com/Revachol/iu5_web_5sem/internal/app/ds"

func (r *Repository) GetServicesByRequestID(requestID int) ([]ds.Historical_request_service, error) {
	var services []ds.Historical_request_service
	err := r.db.Preload("Service").Where("request_id = ?", requestID).Find(&services).Error
	return services, err
}
