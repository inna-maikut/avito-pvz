//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestNewProductRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := NewProductRepository(&sqlx.DB{}, &trmsqlx.CtxGetter{})
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("error.first_nil", func(t *testing.T) {
		res, err := NewProductRepository(nil, &trmsqlx.CtxGetter{})
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.second_nil", func(t *testing.T) {
		res, err := NewProductRepository(&sqlx.DB{}, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestProductRepository_Create(t *testing.T) {
	db := setUp(t)
	repo, err := NewProductRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	ID1 := model.NewPVZID()
	receptionID1 := model.NewReceptionID()
	receptionID2 := model.NewReceptionID()

	type args struct {
		receptionID model.ReceptionID
		category    model.ProductCategory
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		check   func(t *testing.T, res model.Product)
		wantErr bool
	}{
		{
			name: "success_insert",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`INSERT INTO pvz(id, city) VALUES($1, $2)`, ID1, "Москва")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`,
					receptionID1, ID1, model.ReceptionStatusInProgress)
				require.NoError(t, err)
			},
			args: args{
				receptionID: receptionID1,
				category:    model.ProductCategoryElectronics,
			},
			check: func(t *testing.T, res model.Product) {
				var product Product
				err = db.Get(&product, "SELECT id, reception_id, category, added_at FROM products WHERE reception_id = $1", receptionID1)
				require.NoError(t, err)

				require.Equal(t, receptionID1.UUID(), product.ReceptionID)
				require.Equal(t, int16(model.ProductCategoryElectronics), product.Category)
				require.Equal(t, product.ID, res.ID.UUID())
				require.Equal(t, receptionID1, res.ReceptionID)
				require.Equal(t, model.ProductCategoryElectronics, res.Category)
			},
			wantErr: false,
		},
		{
			name: "no_reception_error",
			prepare: func(_ *testing.T) {
			},
			args: args{
				receptionID: receptionID2,
				category:    model.ProductCategoryElectronics,
			},
			check: func(_ *testing.T, _ model.Product) {
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.Create(context.Background(), tc.args.receptionID, tc.args.category)

			require.Equal(t, err != nil, tc.wantErr)
			tc.check(t, res)
		})
	}
}

func TestProductRepository_RemoveLast(t *testing.T) {
	db := setUp(t)
	repo, err := NewProductRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	ID1 := model.NewPVZID()
	receptionID1 := model.NewReceptionID()
	receptionID2 := model.NewReceptionID()

	now := time.Now().Truncate(time.Second)

	type args struct {
		receptionID model.ReceptionID
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		check   func(t *testing.T)
		wantErr error
	}{
		{
			name: "success_delete",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`INSERT INTO pvz(id, city) VALUES($1, $2)`, ID1, "Москва")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`,
					receptionID1, ID1, model.ReceptionStatusInProgress)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO products(reception_id, category, added_at) VALUES($1, $2, $3)`,
					receptionID1, model.ProductCategoryClothes, now.Add(-time.Second))
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO products(reception_id, category, added_at) VALUES($1, $2, $3)`,
					receptionID1, model.ProductCategoryElectronics, now)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO products(reception_id, category, added_at) VALUES($1, $2, $3)`,
					receptionID1, model.ProductCategoryShoes, now.Add(-2*time.Second))
				require.NoError(t, err)
			},
			args: args{
				receptionID: receptionID1,
			},
			check: func(t *testing.T) {
				var categories []model.ProductCategory
				err = db.Select(&categories, "SELECT category FROM products WHERE reception_id = $1", receptionID1)
				require.NoError(t, err)

				require.ElementsMatch(t, []model.ProductCategory{
					model.ProductCategoryShoes,
					model.ProductCategoryClothes,
				}, categories)
			},
			wantErr: nil,
		},
		{
			name: "no_product_error",
			prepare: func(_ *testing.T) {
			},
			args: args{
				receptionID: receptionID2,
			},
			check: func(_ *testing.T) {
			},
			wantErr: model.ErrProductNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			err = repo.RemoveLast(context.Background(), tc.args.receptionID)

			require.ErrorIs(t, err, tc.wantErr)
			tc.check(t)
		})
	}
}

func TestProductRepository_GetByReceptionIDs(t *testing.T) {
	db := setUp(t)
	repo, err := NewProductRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)
	ID1 := model.NewPVZID()
	receptionID1 := model.NewReceptionID()
	receptionID2 := model.NewReceptionID()
	receptionID3 := model.NewReceptionID()
	productID1 := model.NewProductID()
	productID2 := model.NewProductID()
	productID3 := model.NewProductID()
	productID4 := model.NewProductID()

	now := time.Now().Truncate(time.Second)

	type args struct {
		receptionIDs []model.ReceptionID
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		wantErr error
		wantRes []model.Product
	}{
		{
			name: "success",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`INSERT INTO pvz(id, city) VALUES($1, $2)`, ID1, "Москва")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`,
					receptionID1, ID1, model.ReceptionStatusInProgress)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`,
					receptionID2, ID1, model.ReceptionStatusInProgress)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO receptions(id, pvz_id, status) VALUES($1, $2, $3)`,
					receptionID3, ID1, model.ReceptionStatusInProgress)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO products(id, reception_id, category, added_at) VALUES($1, $2, $3, $4)`,
					productID1, receptionID1, model.ProductCategoryClothes, now.Add(-time.Second))
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO products(id, reception_id, category, added_at) VALUES($1, $2, $3, $4)`,
					productID2, receptionID1, model.ProductCategoryElectronics, now)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO products(id, reception_id, category, added_at) VALUES($1, $2, $3, $4)`,
					productID3, receptionID2, model.ProductCategoryShoes, now.Add(2*time.Second))
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO products(id, reception_id, category, added_at) VALUES($1, $2, $3, $4)`,
					productID4, receptionID3, model.ProductCategoryShoes, now.Add(2*time.Second))
				require.NoError(t, err)
			},
			args: args{
				receptionIDs: []model.ReceptionID{
					receptionID1, receptionID2,
				},
			},
			wantErr: nil,
			wantRes: []model.Product{
				{
					ID:          productID1,
					ReceptionID: receptionID1,
					Category:    model.ProductCategoryClothes,
					AddedAt:     now.Add(-time.Second),
				},
				{
					ID:          productID2,
					ReceptionID: receptionID1,
					Category:    model.ProductCategoryElectronics,
					AddedAt:     now,
				},
				{
					ID:          productID3,
					ReceptionID: receptionID2,
					Category:    model.ProductCategoryShoes,
					AddedAt:     now.Add(2 * time.Second),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.GetByReceptionIDs(context.Background(), tc.args.receptionIDs)

			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.wantRes, res)
		})
	}
}
