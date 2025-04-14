package repository

import (
	"context"
	"errors"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type ProductRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewProductRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*ProductRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &ProductRepository{
		db:     db,
		getter: getter,
	}, nil
}

func (r *ProductRepository) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}

func (r *ProductRepository) Create(ctx context.Context, receptionID model.ReceptionID, category model.ProductCategory) (model.Product, error) {
	var product Product

	q := `INSERT INTO products (reception_id, category) VALUES ($1, $2)
	RETURNING id, reception_id, category, added_at`

	err := r.trOrDB(ctx).GetContext(ctx, &product, q, receptionID, category)
	if err != nil {
		return model.Product{}, fmt.Errorf("db.GetContext: %w", err)
	}

	return model.Product{
		ID:          model.ProductID(product.ID),
		ReceptionID: receptionID,
		Category:    category,
		AddedAt:     product.AddedAt,
	}, nil
}

func (r *ProductRepository) RemoveLast(ctx context.Context, receptionID model.ReceptionID) error {
	q := `DELETE FROM products WHERE id IN (
		SELECT id FROM products WHERE reception_id = $1 ORDER BY added_at DESC LIMIT 1
	)`

	result, err := r.trOrDB(ctx).ExecContext(ctx, q, receptionID)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}
	if count == 0 {
		return model.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) GetByReceptionIDs(ctx context.Context, receptionIDs []model.ReceptionID) ([]model.Product, error) {
	var entities []Product

	q := "SELECT id, reception_id, category, added_at FROM products WHERE reception_id = ANY($1::UUID[]) ORDER BY added_at"

	err := r.trOrDB(ctx).SelectContext(ctx, &entities, q, receptionIDs)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	products := make([]model.Product, 0, len(entities))

	for _, product := range entities {
		products = append(products, model.Product{
			ID:          model.ProductID(product.ID),
			ReceptionID: model.ReceptionID(product.ReceptionID),
			Category:    model.ProductCategory(product.Category),
			AddedAt:     product.AddedAt,
		})
	}
	return products, nil
}
