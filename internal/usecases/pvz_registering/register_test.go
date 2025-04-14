package pvz_registering

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestUseCase_RegisterPVZ(t *testing.T) {
	type mocks struct {
		pvzRepo *MockpvzRepo
	}
	type args struct {
		city string
	}

	testCases := []struct {
		name     string
		prepare  func(t *testing.T, m *mocks)
		args     args
		wantErr  error
		wantCity string
	}{
		{
			name: "success.register",
			prepare: func(t *testing.T, m *mocks) {
				m.pvzRepo.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pvz model.PVZ) error {
						require.Equal(t, "test1", pvz.City)
						require.WithinDuration(t, time.Now(), pvz.RegisteredAt, time.Minute)
						return nil
					})
			},
			args: args{
				city: "test1",
			},
			wantErr:  nil,
			wantCity: "test1",
		},
		{
			name: "error.Register",
			prepare: func(_ *testing.T, m *mocks) {
				m.pvzRepo.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			args: args{
				city: "test1",
			},
			wantErr:  assert.AnError,
			wantCity: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				pvzRepo: NewMockpvzRepo(ctrl),
			}

			tc.prepare(t, m)

			uc, err := New(m.pvzRepo)
			require.NoError(t, err)

			pvz, err := uc.RegisterPVZ(context.Background(), tc.args.city)
			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.wantCity, pvz.City)
		})
	}
}
