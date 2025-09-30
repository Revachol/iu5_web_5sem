package repository

import (
	"fmt"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/sirupsen/logrus"
)

func (r *Repository) GetOrders() ([]ds.Historical_service, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	var orders []ds.Historical_service
	err := r.db.Find(&orders).Error
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("Array is empty")
	}
	return orders, nil
}

func (r *Repository) GetOrder(id int) (ds.Historical_service, error) {
	var estimate ds.Historical_service
	err := r.db.Where("id = ?", id).First(&estimate).Error
	if err != nil {
		logrus.Error("Error fetching estimate:", err)
		return ds.Historical_service{}, err
	}
	return estimate, nil
}

func (r *Repository) GetOrdersByTitle(title string) ([]ds.Historical_service, error) {
	var orders []ds.Historical_service
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) GetCartCount() int64 {
	var (
		requestID uint
		count     int64
	)
	creatorID := 1

	// Поиск черновой заявки
	err := r.db.Model(&ds.Historical_request{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("id").First(&requestID).Error
	if err != nil {
		return 0
	}

	// Подсчёт количества исторических объектов в заявке
	err = r.db.Model(&ds.Historical_request_service{}).
		Where("request_id = ?", requestID).
		Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in Historical_service_Historical_service_request:", err)
	}

	return count
}

func (r *Repository) GetDraftRequestID() (int, error) {
	var request ds.Historical_request
	creatorID := 1
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").First(&request).Error
	if err != nil {
		return 0, err
	}
	return request.ID, nil
}
