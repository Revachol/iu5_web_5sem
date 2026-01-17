package ds

type Historical_service struct {
	ID               int     `gorm:"primaryKey;autoIncrement" json:"ID"`
	Name             string  `gorm:"type:varchar(255);not null" json:"NameOfHistoricalObject"`
	Description      string  `gorm:"type:text" json:"DescriptionOfHistoricalObject"`
	PriceUSD         float64 `gorm:"type:decimal(15,2);not null" json:"PriceUSDOfHistoricalObject"`
	Unit             string  `gorm:"type:varchar(100);not null" json:"UnitOfHistoricalObject"`
	HistoricalPeriod string  `gorm:"type:varchar(100);not null" json:"HistoricalPeriodOfHistoricalObject"`
	HistoricalRegion string  `gorm:"type:varchar(100);not null" json:"HistoricalRegionOfHistoricalObject"`
	DataSource       string  `gorm:"type:text" json:"DataSourceOfHistoricalObject"`
	ImageURL         string  `gorm:"type:varchar(500)" json:"ImageURLOfHistoricalObject"`
}
