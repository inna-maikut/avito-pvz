package registering

import (
	"context"
	"errors"
	"fmt"

	"github.com/inna-maikut/avito-pvz/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	userRepo userRepo
}

func New(userRepo userRepo) (*UseCase, error) {
	if userRepo == nil {
		return nil, errors.New("userRepo is nil")
	}
	return &UseCase{
		userRepo: userRepo,
	}, nil
}

func (uc *UseCase) Register(ctx context.Context, email, password string, role model.UserRole) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	user, err := uc.userRepo.Create(ctx, email, string(hashedPassword), role)
	if err != nil {
		return nil, fmt.Errorf("userRepo.Create: %w", err)
	}

	return user, nil
}
