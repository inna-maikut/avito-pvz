//go:build integration

package repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestNewPVZLocker(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := NewPVZLocker(&sqlx.DB{}, &trmsqlx.CtxGetter{})
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("error.first_nil", func(t *testing.T) {
		res, err := NewPVZLocker(nil, &trmsqlx.CtxGetter{})
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.second_nil", func(t *testing.T) {
		res, err := NewPVZLocker(&sqlx.DB{}, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestPVZLocker_Lock(t *testing.T) {
	db := setUp(t)
	locker, err := NewPVZLocker(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	ID1 := model.NewPVZID()

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))

	err = trManager.Do(context.Background(), func(ctx context.Context) error {
		err = locker.Lock(ctx, ID1)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
}
