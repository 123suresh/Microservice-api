package utils

import (
	"errors"
	"os"
	"time"

	"example.com/dynamicWordpressBuilding/internal/model"
	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func NewPayload(id uint, username string, duration time.Duration) (*model.Payload, error) {
	payload := &model.Payload{
		ID:        id,
		Email:     username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (maker *JWTMaker) CreateToken(id uint, username string, duration time.Duration) (string, error) {
	tokenPayload, err := NewPayload(id, username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenPayload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*model.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &model.Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*model.Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

//if you have to use CreateToken function then first make Interface We have make as maker.go
//then initialize interface like below and you will be able to use
func NewTokenMaker() Maker {
	return &JWTMaker{
		secretKey: os.Getenv("SECRET_KEY"),
	}
}
