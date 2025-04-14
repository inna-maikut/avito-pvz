package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type PVZLocker struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewPVZLocker(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*PVZLocker, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &PVZLocker{
		db:     db,
		getter: getter,
	}, nil
}

func (r *PVZLocker) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}

func (r *PVZLocker) Lock(ctx context.Context, pvzID model.PVZID) error {
	q := `SELECT pg_advisory_xact_lock($1)`

	_, err := r.trOrDB(ctx).ExecContext(ctx, q, pvzLockerHash(pvzID))
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

func pvzLockerHash(pvzID model.PVZID) int64 {
	bigIntUUID := new(big.Int).SetBytes(pvzID[:])
	return bigIntUUID.Int64()
}
