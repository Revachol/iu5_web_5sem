package ds

import "time"

type Historical_request struct {
	ID                  int       `gorm:"primaryKey;autoIncrement"`
	Status              string    `gorm:"type:varchar(20);not null;check:status IN ('draft', 'deleted', 'submitted', 'completed', 'rejected')"`
	CreatorID           int       `gorm:"not null"`
	Creator             User      `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt           time.Time `gorm:"not null;default:now()"`
	SubmittedAt         *time.Time
	CompletedAt         *time.Time
	ModeratorID         *int
	Moderator           *User    `gorm:"foreignKey:ModeratorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TotalCostUSD        *float64 `gorm:"type:decimal(15,2)"`
	InflationMultiplier *float64 `gorm:"type:decimal(10,4)"`
}
