package repository

import (
	"fmt"
	"time"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/sirupsen/logrus"
)

func (r *Repository) GetHistoricalObjects() ([]ds.Historical_service, error) {
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
	err := r.db.Where("name ILIKE ?", "%"+title+"%").Find(&orders).Error
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

func (r *Repository) AddServiceToRequest(serviceID int) error {
	var requestID uint
	creatorID := 1
	moderatorID := 2

	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc).Truncate(time.Second)

	err := r.db.Model(&ds.Historical_request{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("id").
		First(&requestID).Error

	if err != nil {
		moderatorIDptr := &moderatorID
		newReq := ds.Historical_request{
			Status:      "draft",
			CreatedAt:   now,
			CreatorID:   creatorID,
			ModeratorID: moderatorIDptr,
		}
		if err := r.db.Create(&newReq).Error; err != nil {
			return fmt.Errorf("не удалось создать черновую заявку: %w", err)
		}
		requestID = uint(newReq.ID)
	}

	var service ds.Historical_service
	if err := r.db.First(&service, serviceID).Error; err != nil {
		return fmt.Errorf("услуга с id %d не найдена: %w", serviceID, err)
	}

	record := ds.Historical_request_service{
		RequestID:     int(requestID),
		ServiceID:     serviceID,
		Quantity:      1,
		UnitPriceUSD:  service.PriceUSD,
		TotalPriceUSD: service.PriceUSD,
	}
	if err := r.db.Create(&record).Error; err != nil {
		return fmt.Errorf("ошибка при добавлении услуги в заявку: %w", err)
	}

	return nil
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
