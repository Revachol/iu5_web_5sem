package ds

type Historical_request_service struct {
	RequestID     int                `gorm:"primaryKey;not null"`
	ServiceID     int                `gorm:"primaryKey;not null"`
	Request       Historical_request `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Service       Historical_service `gorm:"foreignKey:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity      float64            `gorm:"type:decimal(10,2);not null"`
	UnitPriceUSD  float64            `gorm:"type:decimal(15,2);not null"`
	TotalPriceUSD float64            `gorm:"type:decimal(15,2);not null"`
	DisplayOrder  int                `gorm:"not null;default:0"`
}
