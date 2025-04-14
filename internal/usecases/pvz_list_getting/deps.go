//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package pvz_list_getting

import (
	"context"
	"time"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type pvzRepo interface {
	Get(ctx context.Context, pvzIDs []model.PVZID) ([]model.PVZ, error)
}

type receptionRepo interface {
	Search(ctx context.Context, receptedAtFrom, receptedAtTo *time.Time, offset, limit int64) ([]model.Reception, error)
}

type productRepo interface {
	GetByReceptionIDs(ctx context.Context, receptionIDs []model.ReceptionID) ([]model.Product, error)
}
