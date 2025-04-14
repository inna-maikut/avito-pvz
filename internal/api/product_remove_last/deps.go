//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package product_remove_last

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type productRemoving interface {
	RemoveLastProduct(ctx context.Context, pvzID model.PVZID) error
}
