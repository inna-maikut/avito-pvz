package dummy_authenticating

import (
	"context"
	"errors"
	"fmt"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type UseCase struct {
	tokenProvider tokenProvider
}

func New(tokenProvider tokenProvider) (*UseCase, error) {
	if tokenProvider == nil {
		return nil, errors.New("tokenProvider is nil")
	}
	return &UseCase{
		tokenProvider: tokenProvider,
	}, nil
}

func (uc *UseCase) Auth(_ context.Context, role model.UserRole) (string, error) {
	token, err := uc.tokenProvider.CreateToken("", 0, role)
	if err != nil {
		return "", fmt.Errorf("tokenProvider.CreateToken: %w", err)
	}

	return token, nil
}
