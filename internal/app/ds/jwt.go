package ds

import "github.com/golang-jwt/jwt/v5"

// данные сохраняющиеся внутри JWT
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

// ответ после аутентификации
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
