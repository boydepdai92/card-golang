package jwt

import (
	"card-warhouse/components/tokenProvider"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type jwtProvider struct {
	secret string
}

func NewJwtProvider(secret string) *jwtProvider {
	return &jwtProvider{secret: secret}
}

type myClaims struct {
	Payload tokenProvider.TokenPayload
	jwt.StandardClaims
}

func (j *jwtProvider) Generate(data tokenProvider.TokenPayload, expiresIn int) (*tokenProvider.Token, error) {
	createdAt := time.Now()

	tokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		data,
		jwt.StandardClaims{
			ExpiresAt: createdAt.Local().Add(time.Second * time.Duration(expiresIn)).Unix(),
			IssuedAt:  createdAt.Local().Unix(),
		},
	})

	token, err := tokenRaw.SignedString([]byte(j.secret))

	if nil != err {
		return nil, err
	}

	return &tokenProvider.Token{
		Token:     token,
		ExpiresIn: expiresIn,
		CreatedAt: createdAt,
	}, nil
}

func (j *jwtProvider) Validate(token string) (*tokenProvider.TokenPayload, error) {
	result, err := jwt.ParseWithClaims(token, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if nil != err {
		return nil, tokenProvider.ErrTokenInvalid
	}

	if !result.Valid {
		return nil, tokenProvider.ErrTokenInvalid
	}

	claims, ok := result.Claims.(*myClaims)

	if !ok {
		return nil, tokenProvider.ErrTokenInvalid
	}

	return &claims.Payload, nil
}
