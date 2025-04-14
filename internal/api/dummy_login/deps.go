//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package dummy_login

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type authenticating interface {
	Auth(ctx context.Context, role model.UserRole) (string, error)
}
