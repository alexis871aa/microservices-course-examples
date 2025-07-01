package model

import "time"

// User - модель пользователя
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"hashed_password"`
}

// TokenPair - пара токенов
type TokenPair struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

// Claims - кастомные claims для JWT
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}
