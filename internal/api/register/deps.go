//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package register

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type registering interface {
	Register(ctx context.Context, email, password string, role model.UserRole) (*model.User, error)
}
