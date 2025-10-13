package handler

import (
	"fmt"
	//"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GET /api/users/:id
func (h *Handler) GetUserAPI(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректный ID пользователя"))
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"is_moderator": user.IsModerator,
	})
}

// POST /api/users/register
func (h *Handler) RegisterUserAPI(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректные данные: %v", err))
		return
	}

	err := h.Repository.CreateUser(req.Email, req.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка при создании пользователя: %v", err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Пользователь успешно зарегистрирован",
	})
}

// POST /api/users/login
func (h *Handler) LoginUserAPI(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректные данные: %v", err))
		return
	}

	user, err := h.Repository.GetUserByEmail(req.Email)
	if err != nil || user.PasswordHash != req.Password {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("неверный email или пароль"))
		return
	}

	// В реальном приложении здесь должен быть механизм создания сессии или JWT

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Аутентификация успешна",
		"user_id": user.ID,
	})
}

// POST /api/logout
func (h *Handler) LogoutUserAPI(ctx *gin.Context) {
	// В реальном приложении здесь должен быть механизм уничтожения сессии или JWT

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Деавторизация успешна",
	})
}

func (h *Handler) UpdateUserAPI(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректный ID пользователя"))
		return
	}

	var req struct {
		Email    *string `json:"email,omitempty" binding:"omitempty,email"`
		Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректные данные: %v", err))
		return
	}

	err = h.Repository.UpdateUser(userID, req.Email, req.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка при обновлении пользователя: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Пользователь успешно обновлён",
	})
}
