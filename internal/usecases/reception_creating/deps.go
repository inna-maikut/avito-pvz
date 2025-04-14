//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package reception_creating

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type trManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

type receptionRepo interface {
	GetInProgress(ctx context.Context, pvzID model.PVZID) (model.Reception, error)
	Create(ctx context.Context, pvzID model.PVZID, status model.ReceptionStatus) (model.Reception, error)
}

type pvzLocker interface {
	Lock(ctx context.Context, pvzID model.PVZID) error
}
