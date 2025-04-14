package product_adding

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMocktrManager(ctrl), NewMockreceptionRepo(ctrl), NewMockpvzLocker(ctrl), NewMockproductRepo(ctrl))
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("error.first_nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(nil, NewMockreceptionRepo(ctrl), NewMockpvzLocker(ctrl), NewMockproductRepo(ctrl))
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.second_nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMocktrManager(ctrl), nil, NewMockpvzLocker(ctrl), NewMockproductRepo(ctrl))
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.third_nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMocktrManager(ctrl), NewMockreceptionRepo(ctrl), nil, NewMockproductRepo(ctrl))
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.fourth_nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMocktrManager(ctrl), NewMockreceptionRepo(ctrl), NewMockpvzLocker(ctrl), nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestUseCase_AddProduct(t *testing.T) {
	type mocks struct {
		trManager     *MocktrManager
		receptionRepo *MockreceptionRepo
		productRepo   *MockproductRepo
		pvzLocker     *MockpvzLocker
	}
	type args struct {
		pvzID    model.PVZID
		category model.ProductCategory
	}

	ID1 := model.NewPVZID()
	productID := model.NewProductID()
	receptionID1 := model.NewReceptionID()
	now := time.Now()

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantErr error
		wantRes model.Product
	}{
		{
			name: "success",
			prepare: func(m *mocks) {
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.pvzLocker.EXPECT().
					Lock(gomock.Any(), ID1).
					Return(nil)
				m.receptionRepo.EXPECT().
					GetInProgress(gomock.Any(), ID1).
					Return(model.Reception{
						ID:              receptionID1,
						PVZID:           ID1,
						ReceptionStatus: model.ReceptionStatusInProgress,
						ReceptedAt:      now,
					}, nil)
				m.productRepo.EXPECT().
					Create(gomock.Any(), receptionID1, model.ProductCategoryElectronics).
					Return(model.Product{
						ID:          productID,
						ReceptionID: receptionID1,
						Category:    model.ProductCategoryElectronics,
						AddedAt:     now,
					}, nil)
			},
			args: args{
				pvzID:    ID1,
				category: model.ProductCategoryElectronics,
			},
			wantErr: nil,
			wantRes: model.Product{
				ID:          productID,
				ReceptionID: receptionID1,
				Category:    model.ProductCategoryElectronics,
				AddedAt:     now,
			},
		},
		{
			name: "businessError.ErrReceptionNotFound",
			prepare: func(m *mocks) {
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.pvzLocker.EXPECT().
					Lock(gomock.Any(), ID1).
					Return(nil)
				m.receptionRepo.EXPECT().
					GetInProgress(gomock.Any(), ID1).
					Return(model.Reception{}, model.ErrReceptionNotFound)
			},
			args: args{
				pvzID:    ID1,
				category: model.ProductCategoryElectronics,
			},
			wantErr: model.ErrReceptionNotFound,
			wantRes: model.Product{},
		},
		{
			name: "error.Lock",
			prepare: func(m *mocks) {
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.pvzLocker.EXPECT().
					Lock(gomock.Any(), ID1).
					Return(assert.AnError)
			},
			args: args{
				pvzID:    ID1,
				category: model.ProductCategoryElectronics,
			},
			wantErr: assert.AnError,
			wantRes: model.Product{},
		},
		{
			name: "error.GetInProgress",
			prepare: func(m *mocks) {
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.pvzLocker.EXPECT().
					Lock(gomock.Any(), ID1).
					Return(nil)
				m.receptionRepo.EXPECT().
					GetInProgress(gomock.Any(), ID1).
					Return(model.Reception{}, assert.AnError)
			},
			args: args{
				pvzID:    ID1,
				category: model.ProductCategoryElectronics,
			},
			wantErr: assert.AnError,
			wantRes: model.Product{},
		},
		{
			name: "error.Create",
			prepare: func(m *mocks) {
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.pvzLocker.EXPECT().
					Lock(gomock.Any(), ID1).
					Return(nil)
				m.receptionRepo.EXPECT().
					GetInProgress(gomock.Any(), ID1).
					Return(model.Reception{
						ID:              receptionID1,
						PVZID:           ID1,
						ReceptionStatus: model.ReceptionStatusInProgress,
						ReceptedAt:      now,
					}, nil)
				m.productRepo.EXPECT().
					Create(gomock.Any(), receptionID1, model.ProductCategoryElectronics).
					Return(model.Product{}, assert.AnError)
			},
			args: args{
				pvzID:    ID1,
				category: model.ProductCategoryElectronics,
			},
			wantErr: assert.AnError,
			wantRes: model.Product{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				trManager:     NewMocktrManager(ctrl),
				receptionRepo: NewMockreceptionRepo(ctrl),
				productRepo:   NewMockproductRepo(ctrl),
				pvzLocker:     NewMockpvzLocker(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.trManager, m.receptionRepo, m.pvzLocker, m.productRepo)
			require.NoError(t, err)

			product, err := uc.AddProduct(context.Background(), tc.args.pvzID, tc.args.category)
			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.wantRes, product)
		})
	}
}
