//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package product_add

import (
	"context"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type productAdding interface {
	AddProduct(ctx context.Context, pvzID model.PVZID, category model.ProductCategory) (model.Product, error)
}
