package tokenProvider

import (
	"errors"
	"time"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
)

type Provider interface {
	Generate(data TokenPayload, expiresIn int) (*Token, error)
	Validate(token string) (*TokenPayload, error)
}

type TokenPayload struct {
	UserId int `json:"user_id"`
}

type Token struct {
	Token     string    `json:"token"`
	ExpiresIn int       `json:"expires_in"`
	CreatedAt time.Time `json:"created_at"`
}
