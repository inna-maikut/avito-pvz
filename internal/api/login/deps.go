//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package login

import (
	"context"
)

type authenticating interface {
	Auth(ctx context.Context, email, password string) (string, error)
}
