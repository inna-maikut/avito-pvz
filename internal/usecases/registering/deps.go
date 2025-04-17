//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package registering

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type userRepo interface {
	Create(ctx context.Context, email, passwordHash string, role model.UserRole) (*model.User, error)
}
