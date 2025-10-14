package repository

import (
	"errors"
	"fmt"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/Revachol/iu5_web_5sem/internal/app/role"
	"golang.org/x/crypto/bcrypt"
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

func (r *Repository) CreateUser(email, password string, Role role.Role) (*ds.User, error) {
	var existing ds.User
	err := r.db.Where("email = ?", email).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Настоящая ошибка БД
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing.ID != 0 {
		// Пользователь уже существует
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Хешируем пароль перед сохранением
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := ds.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         Role,
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
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

func (r *Repository) Authenticate(email, password string) (*ds.User, error) {
	var user ds.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Сравниваем введённый пароль с хешем
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}
