//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package middleware

import (
	"github.com/inna-maikut/avito-pvz/internal/model"
)

type tokenProvider interface {
	ParseToken(tokenStr string) (model.TokenInfo, error)
}
