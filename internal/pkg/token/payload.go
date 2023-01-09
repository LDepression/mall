package token

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	// 用于管理每个JWT
	ID       uuid.UUID
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	// 创建时间用于检验
	IssuedAt  time.Time `json:"issued-at"`
	ExpiredAt time.Time `json:"expired-at"`
}

func NewPayload(userID int64, userName string, expireDate time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenID,
		UserID:    userID,
		UserName:  userName,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(expireDate),
	}, nil
}
