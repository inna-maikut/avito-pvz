//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package pvz_registering

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type pvzRepo interface {
	Register(ctx context.Context, pvz model.PVZ) error
}
