//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package reception_create

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type receptionCreating interface {
	CreateReception(ctx context.Context, pvzID model.PVZID) (model.Reception, error)
}
