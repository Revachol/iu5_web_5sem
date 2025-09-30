package repository

import (
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

func (r *Repository) GetHistoricalRequest(id int) (ds.Historical_service, error) {
	var historical_object ds.Historical_service
	err := r.db.Where("id = ?", id).First(&historical_object).Error
	if err != nil {
		return ds.Historical_service{}, err
	}
	return historical_object, nil
}
