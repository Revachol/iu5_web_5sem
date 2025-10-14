package ds

import (
	"github.com/Revachol/iu5_web_5sem/internal/app/role"
	"time"
)

type User struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	FullName     string    `gorm:"type:varchar(255)"`
	IsModerator  bool      `gorm:"not null;default:false"`
	IsActive     bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
	Role         role.Role `gorm:"int;not null;default:0" json:"role"`
}
