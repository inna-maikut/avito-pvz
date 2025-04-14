//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package pvz_register

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type pvzRegistering interface {
	RegisterPVZ(ctx context.Context, city string) (model.PVZ, error)
}
