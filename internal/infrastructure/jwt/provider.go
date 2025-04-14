package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

const (
	tokenLifetime = time.Hour * 72
)

var (
	ErrInvalidJWTToken           = errors.New("invalid JWT token")
	ErrInvalidUserIDInJWTToken   = errors.New("invalid userID in JWT token")
	ErrInvalidUsernameInJWTToken = errors.New("invalid username in JWT token")
	ErrInvalidRoleInJWTToken     = errors.New("invalid role in JWT token")
)

type Provider struct {
	secret []byte
}

func NewProviderFromEnv() (*Provider, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("env JWT_SECRET is empty")
	}

	provider := &Provider{
		secret: []byte(secret),
	}

	return provider, nil
}

func (p *Provider) CreateToken(email string, userID int64, role model.UserRole) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(tokenLifetime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(p.secret)
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return tokenStr, nil
}

func (p *Provider) ParseToken(tokenStr string) (model.TokenInfo, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (any, error) {
		return p.secret, nil
	})
	if err != nil {
		return model.TokenInfo{}, fmt.Errorf("jwt.Parse: %w", err)
	}

	if !token.Valid {
		return model.TokenInfo{}, ErrInvalidJWTToken
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	userID, ok := claims["userID"].(float64)
	if !ok {
		return model.TokenInfo{}, ErrInvalidUserIDInJWTToken
	}

	username, ok := claims["email"].(string)
	if !ok {
		return model.TokenInfo{}, ErrInvalidUsernameInJWTToken
	}

	role, ok := claims["role"].(string)
	if !ok {
		return model.TokenInfo{}, ErrInvalidRoleInJWTToken
	}

	return model.TokenInfo{
		UserID:   int64(userID),
		Email:    username,
		UserRole: model.UserRole(role),
	}, nil
}
