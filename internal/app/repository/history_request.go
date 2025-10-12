package repository

import (
	"fmt"
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
