package reply

import (
	"mall/internal/pkg/token"
)

type User struct {
	UserID       int64         `json:"userID"`
	UserName     string        `json:"userName"`
	Avatar       string        `json:"avatar"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	Payload      token.Payload `json:"payload"`
	Birthday     string        `json:"birthday"`
	Email        string        `json:"email"`
}

type UserInfo2Visitor struct {
	UserID   int64  `json:"userID"`
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
	Birthday string `json:"birthday"`
	Email    string `json:"email"`
}

type ReplyToken struct {
	NewToken string `json:"newToken"`
}
