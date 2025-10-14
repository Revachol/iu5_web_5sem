package handler

import (
	"fmt"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/Revachol/iu5_web_5sem/internal/app/role"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

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

	user, err := h.Repository.CreateUser(req.Email, req.Password, role.User)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка при создании пользователя: %v", err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Пользователь успешно зарегистрирован",
		"id":      user.ID,
		"login":   user.Email,
	})
}

// POST /api/users/login
func (h *Handler) LoginUserAPI(ctx *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.Repository.Authenticate(body.Email, body.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.GenerateTokens(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.SetCookie(
		cookieName,
		accessToken,
		int(h.Config.JWT.AccessTokenTTL.Seconds()),
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"email":       user.Email,
		"role":        user.Role,
		"accessToken": accessToken,
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

func (h *Handler) GenerateTokens(userID int, role role.Role) (string, error) {
	now := time.Now()

	claims := ds.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(h.Config.JWT.AccessTokenTTL)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.Config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
