package authenticating

import (
	"context"
	"errors"
	"fmt"

	"github.com/inna-maikut/avito-pvz/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	userRepo      userRepo
	tokenProvider tokenProvider
}

func New(userRepo userRepo, tokenProvider tokenProvider) (*UseCase, error) {
	if userRepo == nil {
		return nil, errors.New("userRepo is nil")
	}
	if tokenProvider == nil {
		return nil, errors.New("tokenProvider is nil")
	}
	return &UseCase{
		userRepo:      userRepo,
		tokenProvider: tokenProvider,
	}, nil
}

func (uc *UseCase) Auth(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("userRepo.GetByEmail: %w", err)
	}

	err = uc.checkUserPassword(user.Password, password)
	if err != nil {
		return "", fmt.Errorf("checkUserPassword: %w", err)
	}

	token, err := uc.tokenProvider.CreateToken(user.Email, user.UserID, user.UserRole)
	if err != nil {
		return "", fmt.Errorf("tokenProvider.CreateToken: %w", err)
	}

	return token, nil
}

func (uc *UseCase) checkUserPassword(dbPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return model.ErrWrongUserPassword
		}

		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
	}

	return nil
}
