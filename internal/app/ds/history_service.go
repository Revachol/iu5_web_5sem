package ds

import "time"

type Historical_service struct {
	ID               int       `gorm:"primaryKey;autoIncrement"`
	Name             string    `gorm:"type:varchar(255);not null"`
	Description      string    `gorm:"type:text"`
	PriceUSD         float64   `gorm:"type:decimal(15,2);not null"`
	Unit             string    `gorm:"type:varchar(100);not null"`
	HistoricalPeriod string    `gorm:"type:varchar(100);not null"`
	HistoricalRegion string    `gorm:"type:varchar(100);not null"`
	DataSource       string    `gorm:"type:text"`
	ImageURL         string    `gorm:"type:varchar(500)"`
	IsActive         bool      `gorm:"not null;default:true"`
	CreatedAt        time.Time `gorm:"not null;default:now()"`
}
