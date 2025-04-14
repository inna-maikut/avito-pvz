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

func TestPVZRepository_Register(t *testing.T) {
	db := setUp(t)
	repo, err := NewPVZRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	ID1, err := model.ParsePVZID("0cd22cf8-2636-47ac-9c06-ca0a3e11a19c")
	require.NoError(t, err)

	now := time.Now().Truncate(time.Second)

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		pvz     model.PVZ
		check   func(t *testing.T)
		wantErr error
	}{
		{
			name: "success_insert",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM pvz where id = $1`, ID1)
				require.NoError(t, err)
			},
			pvz: model.PVZ{
				ID:           ID1,
				City:         "test1",
				RegisteredAt: now,
			},
			check: func(t *testing.T) {
				var pvz PVZ
				err = db.Get(&pvz, "SELECT id, city, registered_at FROM pvz WHERE id = $1", ID1)
				require.NoError(t, err)

				require.Equal(t, PVZ{
					ID:           ID1.UUID(),
					City:         "test1",
					RegisteredAt: now,
				}, pvz)
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			err := repo.Register(context.Background(), tc.pvz)

			require.ErrorIs(t, err, tc.wantErr)
			tc.check(t)
		})
	}
}

func TestPVZRepository_Get(t *testing.T) {
	db := setUp(t)
	repo, err := NewPVZRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	pvzID1 := model.NewPVZID()
	pvzID2 := model.NewPVZID()
	pvzID3 := model.NewPVZID()

	now := time.Now().Truncate(time.Second)

	type args struct {
		pvzIDs []model.PVZID
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		wantErr error
		wantRes []model.PVZ
	}{
		{
			name: "success",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`INSERT INTO pvz(id, city, registered_at) VALUES($1, $2, $3)`, pvzID1, "Москва", now)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO pvz(id, city, registered_at) VALUES($1, $2, $3)`, pvzID2, "Москва", now.Add(time.Minute))
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO pvz(id, city, registered_at) VALUES($1, $2, $3)`, pvzID3, "Москва", now.Add(time.Minute*2))
				require.NoError(t, err)
			},
			args: args{
				pvzIDs: []model.PVZID{
					pvzID1, pvzID2,
				},
			},
			wantErr: nil,
			wantRes: []model.PVZ{
				{
					ID:           pvzID1,
					City:         "Москва",
					RegisteredAt: now,
				},
				{
					ID:           pvzID2,
					City:         "Москва",
					RegisteredAt: now.Add(time.Minute),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.Get(context.Background(), tc.args.pvzIDs)

			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.wantRes, res)
		})
	}
}
