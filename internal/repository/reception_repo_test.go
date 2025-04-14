//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestReceptionRepository_GetInProgress(t *testing.T) {
	db := setUp(t)
	repo, err := NewReceptionRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	ID1 := model.NewPVZID()
	receptionID1 := model.NewReceptionID()
	receptionID2 := model.NewReceptionID()

	type args struct {
		pvzID model.PVZID
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		check   func(t *testing.T, res model.Reception)
		wantErr error
	}{
		{
			name: "success",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM receptions where pvz_id = $1`, ID1)
				require.NoError(t, err)
				_, err = db.Exec(`DELETE FROM pvz where id = $1`, ID1)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO pvz(id, city) VALUES($1, $2)`, ID1, "Москва")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`, receptionID1, ID1, model.ReceptionStatusInProgress)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`, receptionID2, ID1, model.ReceptionStatusClose)
				require.NoError(t, err)
			},
			args: args{
				pvzID: ID1,
			},
			check: func(t *testing.T, res model.Reception) {
				require.Equal(t, receptionID1, res.ID)
				require.Equal(t, ID1, res.PVZID)
				require.Equal(t, model.ReceptionStatusInProgress, res.ReceptionStatus)
			},
			wantErr: nil,
		},
		{
			name: "error.notFound",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM receptions where pvz_id = $1`, ID1)
				require.NoError(t, err)
				_, err = db.Exec(`DELETE FROM pvz where id = $1`, ID1)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO pvz(id, city) VALUES($1, $2)`, ID1, "Москва")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`, receptionID2, ID1, model.ReceptionStatusClose)
				require.NoError(t, err)
			},
			args: args{
				pvzID: ID1,
			},
			check: func(_ *testing.T, _ model.Reception) {
			},
			wantErr: model.ErrReceptionNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.GetInProgress(context.Background(), tc.args.pvzID)

			require.ErrorIs(t, err, tc.wantErr)
			tc.check(t, res)
		})
	}
}

func TestReceptionRepository_Create(t *testing.T) {
	db := setUp(t)
	repo, err := NewReceptionRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	ID1, err := model.ParsePVZID("0cd22cf8-2636-47ac-9c06-ca0a3e11a18c")
	require.NoError(t, err)

	type args struct {
		pvzID  model.PVZID
		status model.ReceptionStatus
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		check   func(t *testing.T, res model.Reception)
		wantErr bool
	}{
		{
			name: "success_insert",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM receptions where pvz_id = $1`, ID1)
				require.NoError(t, err)
				_, err = db.Exec(`DELETE FROM pvz where id = $1`, ID1)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO pvz(id, city) VALUES($1, $2)`, ID1, "Москва")
				require.NoError(t, err)
			},
			args: args{
				pvzID:  ID1,
				status: model.ReceptionStatusInProgress,
			},
			check: func(t *testing.T, res model.Reception) {
				var reception Reception
				err = db.Get(&reception, "SELECT id, pvz_id, status, recepted_at FROM receptions WHERE pvz_id = $1", ID1)
				require.NoError(t, err)

				require.Equal(t, ID1.UUID(), reception.PVZID)
				require.Equal(t, int64(model.ReceptionStatusInProgress), reception.Status)
				require.Equal(t, reception.ID, res.ID.UUID())
				require.Equal(t, ID1, res.PVZID)
				require.Equal(t, model.ReceptionStatusInProgress, res.ReceptionStatus)
			},
			wantErr: false,
		},
		{
			name: "no_pvz_error",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM receptions where pvz_id = $1`, ID1)
				require.NoError(t, err)
				_, err = db.Exec(`DELETE FROM pvz where id = $1`, ID1)
				require.NoError(t, err)
			},
			args: args{
				pvzID:  ID1,
				status: model.ReceptionStatusClose,
			},
			check: func(_ *testing.T, _ model.Reception) {
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.Create(context.Background(), tc.args.pvzID, tc.args.status)

			require.Equal(t, err != nil, tc.wantErr)
			tc.check(t, res)
		})
	}
}

func TestReceptionRepository_SetStatus(t *testing.T) {
	db := setUp(t)
	repo, err := NewReceptionRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	pvzID := model.NewPVZID()
	receptionID := model.NewReceptionID()

	type args struct {
		receptionID model.ReceptionID
		status      model.ReceptionStatus
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		check   func(t *testing.T)
		wantErr error
	}{
		{
			name: "success_update",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`INSERT INTO pvz(id, city) VALUES($1, $2)`, pvzID, "Москва")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`,
					receptionID, pvzID, model.ReceptionStatusInProgress)
				require.NoError(t, err)
			},
			args: args{
				receptionID: receptionID,
				status:      model.ReceptionStatusClose,
			},
			check: func(t *testing.T) {
				var reception Reception
				err = db.Get(&reception, "SELECT id, pvz_id, status, recepted_at FROM receptions WHERE id = $1", receptionID)
				require.NoError(t, err)

				require.Equal(t, int64(model.ReceptionStatusClose), reception.Status)
			},
			wantErr: nil,
		},
		{
			name: "no_reception_error",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM receptions where id = $1`, receptionID)
				require.NoError(t, err)
			},
			args: args{
				receptionID: receptionID,
				status:      model.ReceptionStatusClose,
			},
			check: func(_ *testing.T) {
			},
			wantErr: model.ErrReceptionNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			err = repo.SetStatus(context.Background(), tc.args.receptionID, tc.args.status)

			require.ErrorIs(t, err, tc.wantErr)

			tc.check(t)
		})
	}
}

func TestReceptionRepository_Search(t *testing.T) {
	db := setUp(t)
	repo, err := NewReceptionRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	pvzID1 := model.NewPVZID()
	receptionID1 := model.NewReceptionID()
	receptionID2 := model.NewReceptionID()
	receptionID3 := model.NewReceptionID()
	receptionID4 := model.NewReceptionID()

	receptedAtFrom := time.Now().Truncate(time.Second)
	receptedAtTo := receptedAtFrom.Add(time.Hour * 24)

	type args struct {
		receptedAtFrom, receptedAtTo *time.Time
		offset, limit                int64
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		wantErr error
		wantRes []model.Reception
	}{
		{
			name: "success.filter_dates",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM products WHERE TRUE`)
				require.NoError(t, err)
				_, err = db.Exec(`DELETE FROM receptions WHERE TRUE`)
				require.NoError(t, err)
				_, err = db.Exec(`DELETE FROM pvz where id = $1`, pvzID1)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO pvz(id, city) VALUES($1, $2)`, pvzID1, "Москва")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status, recepted_at) VALUES($1, $2, $3, $4)`, receptionID1, pvzID1, model.ReceptionStatusInProgress, receptedAtFrom.Add(-time.Hour))
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status, recepted_at) VALUES($1, $2, $3, $4)`, receptionID2, pvzID1, model.ReceptionStatusClose, receptedAtFrom)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status, recepted_at) VALUES($1, $2, $3, $4)`, receptionID3, pvzID1, model.ReceptionStatusClose, receptedAtTo)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status, recepted_at) VALUES($1, $2, $3, $4)`, receptionID4, pvzID1, model.ReceptionStatusClose, receptedAtTo.Add(time.Hour))
				require.NoError(t, err)
			},
			args: args{
				receptedAtFrom: &receptedAtFrom,
				receptedAtTo:   &receptedAtTo,
				offset:         0,
				limit:          30,
			},
			wantErr: nil,
			wantRes: []model.Reception{
				{
					ID:              receptionID2,
					PVZID:           pvzID1,
					ReceptionStatus: model.ReceptionStatusClose,
					ReceptedAt:      receptedAtFrom,
				},
				{
					ID:              receptionID3,
					PVZID:           pvzID1,
					ReceptionStatus: model.ReceptionStatusClose,
					ReceptedAt:      receptedAtTo,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.Search(context.Background(), tc.args.receptedAtFrom, tc.args.receptedAtTo, tc.args.offset, tc.args.limit)

			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.wantRes, res)
		})
	}
}
