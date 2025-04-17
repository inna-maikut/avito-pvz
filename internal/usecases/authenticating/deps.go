//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package authenticating

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type userRepo interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type tokenProvider interface {
	CreateToken(email string, userID model.UserID, role model.UserRole) (string, error)
}
