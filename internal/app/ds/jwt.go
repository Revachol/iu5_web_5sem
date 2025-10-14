package ds

import (
	"github.com/Revachol/iu5_web_5sem/internal/app/role"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims структура для токена
type JWTClaims struct {
	UserID int       `json:"user_id"`
	Role   role.Role `json:"role"`
	jwt.RegisteredClaims
}
