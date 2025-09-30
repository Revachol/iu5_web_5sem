package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// NewRepository создаёт новый репозиторий и подключается к базе данных PostgreSQL
func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключённой базой данных
	return &Repository{
		db: db,
	}, nil
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db: db,
	}, nil
}

type Order struct { // вот наша новая структура
	ID          int    // поля структур, которые передаются в шаблон
	Title       string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	Price       string
	Value       string
	Img         string
	Description string
	Source      string
}

type Estimate struct {
	ID       int
	OrderIDs []int
}

var estimates = []Estimate{
	{
		ID:       1,
		OrderIDs: []int{1, 2},
	},
}

func (r *Repository) GetEstimateData(id int) (Estimate, error) {
	for _, estimate := range estimates {
		if estimate.ID == id {
			return estimate, nil
		}
	}
	return Estimate{}, fmt.Errorf("estimate with id %d not found", id)
}
