package token

import (
	"Wutils/pkg/utils"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)

	// 构造一个userContent
	userID := utils.RandomInt(1, 1000)
	username := utils.RandomOwner()
	userContent := User{
		UserID:   userID,
		UserName: username,
	}
	content, err := json.Marshal(userContent)
	require.NoError(t, err)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	// 构造paseto
	token, _, err := maker.CreateToken(content, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	result := User{}
	err = json.Unmarshal(payload.Content, &result)
	require.NoError(t, err)
	// 检查token的content
	require.NotZero(t, payload.ID)
	require.Equal(t, userContent, result)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestErrPasetoExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)

	// 构造一个userContent
	userID := utils.RandomInt(1, 1000)
	username := utils.RandomOwner()
	userContent := User{
		UserID:   userID,
		UserName: username,
	}
	content, err := json.Marshal(userContent)
	require.NoError(t, err)

	duration := -time.Minute // 负的
	token, _, err := maker.CreateToken(content, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	fmt.Println("token: ", token)
	verifyToken, err := maker.VerifyToken(token)
	require.Error(t, err, ErrInvalidToken)
	require.Nil(t, verifyToken)
}
