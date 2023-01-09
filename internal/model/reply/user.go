package reply

import (
	"mall/internal/pkg/token"
	"time"
)

type User struct {
	UserID       int64         `json:"userID"`
	UserName     string        `json:"userName"`
	Avatar       string        `json:"avatar"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	Payload      token.Payload `json:"payload"`
	Birthday     time.Time     `json:"birthday"`
	Email        string        `json:"email"`
}
