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
