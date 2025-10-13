package repository

import (
	"errors"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"gorm.io/gorm"
)

func (r *Repository) GetUserByID(userID int) (*ds.User, error) {
	var user ds.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}
