package product_removing

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestUseCase_RemoveLastProduct(t *testing.T) {
	type mocks struct {
		trManager     *MocktrManager
		receptionRepo *MockreceptionRepo
		productRepo   *MockproductRepo
		pvzLocker     *MockpvzLocker
	}
	type args struct {
		pvzID model.PVZID
	}

	ID1 := model.NewPVZID()
	receptionID1 := model.NewReceptionID()
	now := time.Now()

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantErr error
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
					RemoveLast(gomock.Any(), receptionID1).
					Return(nil)
			},
			args: args{
				pvzID: ID1,
			},
			wantErr: nil,
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
				pvzID: ID1,
			},
			wantErr: model.ErrReceptionNotFound,
		},
		{
			name: "businessError.ErrProductNotFound",
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
					RemoveLast(gomock.Any(), receptionID1).
					Return(model.ErrProductNotFound)
			},
			args: args{
				pvzID: ID1,
			},
			wantErr: model.ErrProductNotFound,
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
				pvzID: ID1,
			},
			wantErr: assert.AnError,
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
				pvzID: ID1,
			},
			wantErr: assert.AnError,
		},
		{
			name: "error.RemoveLast",
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
					RemoveLast(gomock.Any(), receptionID1).
					Return(assert.AnError)
			},
			args: args{
				pvzID: ID1,
			},
			wantErr: assert.AnError,
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

			err = uc.RemoveLastProduct(context.Background(), tc.args.pvzID)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
