package token

import (
	"Wutils/pkg/utils"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(minSecretKeySize))
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

	// 构造一个token
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

func TestErrExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(minSecretKeySize))
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

	verifyToken, err := maker.VerifyToken(token)
	require.Error(t, err, ErrInvalidToken)
	require.Nil(t, verifyToken)
}

func TestErrInvalidToken(t *testing.T) {
	payload, err := NewPayload([]byte(utils.RandomOwner()), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(utils.RandomString(minSecretKeySize))
	require.NoError(t, err)

	verifyToken, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Equal(t, err, ErrInvalidToken)
	require.Nil(t, verifyToken)
}
