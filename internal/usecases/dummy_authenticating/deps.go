//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package dummy_authenticating

import (
	"github.com/inna-maikut/avito-pvz/internal/model"
)

type tokenProvider interface {
	CreateToken(email string, userID model.UserID, role model.UserRole) (string, error)
}
