package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific content and duration
func (maker *JWTMaker) CreateToken(content []byte, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(content, duration)
	if err != nil {
		return "", nil, err
	}
	// generate a jwt
	// HS256：一种对称加密算法，使用同一个密钥对signature进行加密解密
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// 令牌签名，对令牌的头部和负荷进行签名，保证令牌不会被伪造或篡改
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, nil
}

// VerifyToken 解析和验证 JWT 令牌，并提取其中的声明
// 如果令牌解析过程中发生错误，会根据错误类型进行相应的处理，如判断是否过期错误。
// 最后，如果解析成功，则返回解析后的 Payload 结构体指针；否则，返回相应的错误。
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	//	验证header中的签名是否合法，提供验证令牌所需的密钥
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC) // 断言token.Method 是否为*jwt.SigningMethodHMAC类型
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	// keyFunc函数需要自己实现

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
