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
	CurrentYear         int      `gorm:"not null;default:0"`
}

type OrderResponse struct {
	ID         int        `json:"idHO"`
	Status     string     `json:"status"`
	DateCreate time.Time  `json:"date_create"`
	DateForm   *time.Time `json:"date_form,omitempty"`
	DateFinish *time.Time `json:"date_finish,omitempty"`
}

type UpdateOrderRequest struct {
	CurrentYear *float64 `json:"current_year,omitempty"`
}

type CompleteOrderRequest struct {
	Status      string `json:"status"`
	ModeratorID int    `json:"moderator_id"`
}
