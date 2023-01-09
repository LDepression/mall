package token

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker([]byte(RandomString(32)))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := RandomString(10)
	duration := time.Minute
	userID := RandomInt(1, 1000)

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, _, err := maker.CreateToken(userID, username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.UserName, username)
	require.Equal(t, payload.UserID, userID)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Millisecond)

	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestMaker(t *testing.T) {
	maker, err := NewPasetoMaker([]byte(RandomString(32)))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	username := RandomOwner()
	duration := time.Second
	userID := RandomInt(1, 1000)
	token, _, err := maker.CreateToken(userID, username, duration)
	require.NoError(t, err)
	result, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	time.Sleep(duration * 2)
	result2, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Empty(t, result2)
}

const alphabetic = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 返回min到max之间的一个随机数
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString 生成一个长度为n的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabetic)
	for i := 0; i < n; i++ {
		c := alphabetic[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}
