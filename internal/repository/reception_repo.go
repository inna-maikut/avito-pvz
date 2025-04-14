package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type ReceptionRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewReceptionRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*ReceptionRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &ReceptionRepository{
		db:     db,
		getter: getter,
	}, nil
}

func (r *ReceptionRepository) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}

func (r *ReceptionRepository) GetInProgress(ctx context.Context, pvzID model.PVZID) (model.Reception, error) {
	var reception Reception

	q := `SELECT id, pvz_id, status, recepted_at	
	FROM receptions	
	WHERE pvz_id = $1 AND status = $2 
	LIMIT 1`

	err := r.trOrDB(ctx).GetContext(ctx, &reception, q, pvzID, model.ReceptionStatusInProgress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Reception{}, model.ErrReceptionNotFound
		}
		return model.Reception{}, fmt.Errorf("db.GetContext: %w", err)
	}

	return model.Reception{
		ID:              model.ReceptionID(reception.ID),
		PVZID:           model.PVZID(reception.PVZID),
		ReceptionStatus: model.ReceptionStatus(reception.Status),
		ReceptedAt:      reception.ReceptedAt,
	}, nil
}

func (r *ReceptionRepository) Create(ctx context.Context, pvzID model.PVZID, status model.ReceptionStatus) (model.Reception, error) {
	var reception Reception

	q := `INSERT INTO receptions (pvz_id, status) VALUES ($1, $2)
	RETURNING id, pvz_id, status, recepted_at`

	err := r.trOrDB(ctx).GetContext(ctx, &reception, q, pvzID, status)
	if err != nil {
		return model.Reception{}, fmt.Errorf("db.GetContext: %w", err)
	}

	return model.Reception{
		ID:              model.ReceptionID(reception.ID),
		PVZID:           model.PVZID(reception.PVZID),
		ReceptionStatus: model.ReceptionStatus(reception.Status),
		ReceptedAt:      reception.ReceptedAt,
	}, nil
}

func (r *ReceptionRepository) SetStatus(ctx context.Context, receptionID model.ReceptionID, status model.ReceptionStatus) error {
	q := `UPDATE receptions SET status = $1 WHERE id = $2`

	result, err := r.trOrDB(ctx).ExecContext(ctx, q, status, receptionID)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}
	if count != 1 {
		return model.ErrReceptionNotFound
	}

	return nil
}

func (r *ReceptionRepository) Search(ctx context.Context, receptedAtFrom, receptedAtTo *time.Time, offset, limit int64) ([]model.Reception, error) {
	if offset < 0 {
		return nil, errors.New("offset can't be negative")
	}
	if limit < 1 {
		return nil, errors.New("limit should be positive")
	}
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "pvz_id", "status", "recepted_at").
		From("receptions").
		OrderBy("recepted_at").
		Offset(uint64(offset)).
		Limit(uint64(limit))

	if receptedAtFrom != nil {
		b = b.Where(sq.GtOrEq{
			"recepted_at": *receptedAtFrom,
		})
	}

	if receptedAtTo != nil {
		b = b.Where(sq.LtOrEq{
			"recepted_at": *receptedAtTo,
		})
	}
	q, args, err := b.ToSql()
	if err != nil {
		return nil, fmt.Errorf("b.ToSql: %w", err)
	}

	var entities []Reception
	err = r.trOrDB(ctx).SelectContext(ctx, &entities, q, args...)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	receptions := make([]model.Reception, 0, len(entities))

	for _, reception := range entities {
		receptions = append(receptions, model.Reception{
			ID:              model.ReceptionID(reception.ID),
			PVZID:           model.PVZID(reception.PVZID),
			ReceptionStatus: model.ReceptionStatus(reception.Status),
			ReceptedAt:      reception.ReceptedAt,
		})
	}
	return receptions, nil
}
