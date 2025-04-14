package dummy_authenticating

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMocktokenProvider(ctrl))
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("error.first_nil", func(t *testing.T) {
		res, err := New(nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestUseCase_Auth(t *testing.T) {
	type mocks struct {
		tokenProvider *MocktokenProvider
	}
	type args struct {
		role model.UserRole
	}

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantRes string
		wantErr error
	}{
		{
			name: "success.moderator",
			prepare: func(m *mocks) {
				m.tokenProvider.EXPECT().CreateToken("", int64(0), model.UserRoleModerator).Return("654321", nil)
			},
			args: args{
				role: model.UserRoleModerator,
			},
			wantRes: "654321",
			wantErr: nil,
		},
		{
			name: "success.employee",
			prepare: func(m *mocks) {
				m.tokenProvider.EXPECT().CreateToken("", int64(0), model.UserRoleEmployee).Return("654321", nil)
			},
			args: args{
				role: model.UserRoleEmployee,
			},
			wantRes: "654321",
			wantErr: nil,
		},
		{
			name: "error.token_provider.create_token",
			prepare: func(m *mocks) {
				m.tokenProvider.EXPECT().CreateToken("", int64(0), model.UserRoleEmployee).Return("", assert.AnError)
			},
			args: args{
				role: model.UserRoleEmployee,
			},
			wantRes: "",
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				tokenProvider: NewMocktokenProvider(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.tokenProvider)
			require.NoError(t, err)

			res, err := uc.Auth(context.Background(), tc.args.role)
			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}
