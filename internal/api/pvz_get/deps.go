//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package pvz_get

import (
	"context"
	"time"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type pvzListGetting interface {
	GetPVZList(ctx context.Context, receptedAtFrom, receptedAtTo *time.Time, page, limit int64) (model.PVZList, error)
}
