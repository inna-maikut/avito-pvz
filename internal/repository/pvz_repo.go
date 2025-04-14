package repository

import (
	"context"
	"errors"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type PVZRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewPVZRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*PVZRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &PVZRepository{
		db:     db,
		getter: getter,
	}, nil
}

func (r *PVZRepository) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}

func (r *PVZRepository) Register(ctx context.Context, pvz model.PVZ) error {
	q := `INSERT INTO pvz (id, city, registered_at)
		VALUES ($1, $2, $3)`

	_, err := r.trOrDB(ctx).ExecContext(ctx, q, pvz.ID, pvz.City, pvz.RegisteredAt)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

func (r *PVZRepository) Get(ctx context.Context, pvzIDs []model.PVZID) ([]model.PVZ, error) {
	var entities []PVZ

	q := "SELECT id, city, registered_at FROM pvz WHERE id = ANY($1::UUID[]) ORDER BY registered_at"

	err := r.trOrDB(ctx).SelectContext(ctx, &entities, q, pvzIDs)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	pvzs := make([]model.PVZ, 0, len(entities))

	for _, pvz := range entities {
		pvzs = append(pvzs, model.PVZ{
			ID:           model.PVZID(pvz.ID),
			City:         pvz.City,
			RegisteredAt: pvz.RegisteredAt,
		})
	}
	return pvzs, nil
}
