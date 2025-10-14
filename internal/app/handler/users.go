package handler

import (
	"fmt"
	"strings"

	"net/http"
	"strconv"
	"time"

	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/Revachol/iu5_web_5sem/internal/app/role"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// RegisterUser регистрирует нового пользователя, выдаёт JWT токен и устанавливает его в cookie.
//
// @Summary Регистрация нового пользователя
// @Description Создаёт пользователя с заданной ролью, возвращает JWT access-токен и устанавливает его в cookie.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body object true "Данные для регистрации пользователя"
// @Success 200 {object} object "Пользователь успешно зарегистрирован"
// @Failure 400 {object} object "Некорректный запрос или пользователь уже существует"
// @Failure 500 {object} object "Ошибка на стороне сервера"
// @Router /api/users/register [post]
func (h *Handler) RegisterUserAPI(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     int    `json:"role"` // "user" или "moderator"
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректные данные: %v", err))
		return
	}

	user, err := h.Repository.CreateUser(req.Email, req.Password, role.Role(req.Role))
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

// Login выполняет аутентификацию пользователя по логину и паролю,
// создаёт JWT access-токен и устанавливает его в заголовок.
//
// @Summary Аутентификация пользователя
// @Description Проверяет логин и пароль, возвращает JWT токен и устанавливает его в заголовок для последующих запросов.
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body object true "Учетные данные пользователя"
// @Success 200 {object} object "Успешная аутентификация"
// @Failure 400 {object} object "Некорректный запрос"
// @Failure 401 {object} object "Неверный логин или пароль"
// @Failure 500 {object} object "Ошибка при генерации токена"
// @Router /api/users/login [post]
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

	ctx.Header("Authorization", "Bearer "+accessToken)
	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

// Logout завершает сессию пользователя, добавляя токен в блеклист Redis.
//
// @Summary Завершение сессии пользователя
// @Description Завершает авторизацию: добавляет JWT в блеклист Redis, аннулируя токен.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Успешный выход"
// @Failure 401 {object} map[string]string "Токен отсутствует или недействителен"
// @Failure 500 {object} map[string]string "Ошибка при работе с Redis"
// @Router /api/users/logout [post]
// @Security ApiKeyAuth
func (h *Handler) LogoutUserAPI(ctx *gin.Context) {
	var tokenStr string

	// Если нет в cookie, пробуем из Authorization header
	if tokenStr == "" {
		authHeader := ctx.GetHeader("Authorization")
		if after, ok := strings.CutPrefix(authHeader, jwtPrefix); ok {
			tokenStr = after
		}
	}

	if tokenStr == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}

	// Проверяем валидность токена
	if _, err := jwt.ParseWithClaims(tokenStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.SecretKey), nil
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Добавляем токен в блеклист Redis
	if h.Redis != nil {
		if err := h.Redis.WriteJWTToBlacklist(ctx.Request.Context(), tokenStr, h.Config.JWT.AccessTokenTTL); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to blacklist token"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
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
