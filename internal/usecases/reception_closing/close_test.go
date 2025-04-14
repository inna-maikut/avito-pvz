package reception_closing

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
		res, err := New(NewMocktrManager(ctrl), NewMockreceptionRepo(ctrl), NewMockpvzLocker(ctrl))
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("error.first_nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(nil, NewMockreceptionRepo(ctrl), NewMockpvzLocker(ctrl))
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.second_nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMocktrManager(ctrl), nil, NewMockpvzLocker(ctrl))
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.third_nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMocktrManager(ctrl), NewMockreceptionRepo(ctrl), nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestUseCase_CloseReception(t *testing.T) {
	type mocks struct {
		trManager     *MocktrManager
		receptionRepo *MockreceptionRepo
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
		wantRes model.Reception
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
				m.receptionRepo.EXPECT().
					SetStatus(gomock.Any(), receptionID1, model.ReceptionStatusClose).
					Return(nil)
			},
			args: args{
				pvzID: ID1,
			},
			wantErr: nil,

			wantRes: model.Reception{
				ID:              receptionID1,
				PVZID:           ID1,
				ReceptionStatus: model.ReceptionStatusClose,
				ReceptedAt:      now,
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
				pvzID: ID1,
			},
			wantErr: model.ErrReceptionNotFound,
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
			name: "error.SetStatus",
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
					}, nil)
				m.receptionRepo.EXPECT().
					SetStatus(gomock.Any(), receptionID1, model.ReceptionStatusClose).
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
				pvzLocker:     NewMockpvzLocker(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.trManager, m.receptionRepo, m.pvzLocker)
			require.NoError(t, err)

			reception, err := uc.CloseReception(context.Background(), tc.args.pvzID)
			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.wantRes, reception)
		})
	}
}
