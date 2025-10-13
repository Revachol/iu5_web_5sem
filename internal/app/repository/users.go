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

func (r *Repository) CreateUser(email, password string) error {
	user := ds.User{
		Email:        email,
		PasswordHash: password, // В реальном приложении пароль должен быть захеширован
		IsModerator:  false,
	}
	return r.db.Create(&user).Error
}

func (r *Repository) GetUserByEmail(email string) (*ds.User, error) {
	var user ds.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(id int, email, password *string) error {
	updates := map[string]interface{}{}
	if email != nil {
		updates["email"] = email
	}
	if password != nil {
		updates["password_hash"] = password // В реальном приложении пароль должен быть захеширован
	}
	return r.db.Model(&ds.User{}).Where("id = ?", id).Updates(updates).Error
}
