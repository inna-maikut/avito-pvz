package authenticating

import (
	"context"
	"testing"

	"github.com/inna-maikut/avito-pvz/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUseCase_Auth(t *testing.T) {
	type mocks struct {
		userRepo      *MockuserRepo
		tokenProvider *MocktokenProvider
	}
	type args struct {
		email    string
		password string
	}

	userID1 := model.NewUserID()

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantRes string
		wantErr error
	}{
		{
			name: "success",
			prepare: func(m *mocks) {
				m.userRepo.EXPECT().
					GetByEmail(gomock.Any(), "test1").
					Return(&model.User{
						UserID:   userID1,
						Email:    "test1",
						Password: makePasswordHash("password1"),
						UserRole: model.UserRoleEmployee,
					}, nil)
				m.tokenProvider.EXPECT().CreateToken("test1", userID1, model.UserRoleEmployee).Return("654321", nil)
			},
			args: args{
				email:    "test1",
				password: "password1",
			},
			wantRes: "654321",
			wantErr: nil,
		},
		{
			name: "error.UserNotFound",
			prepare: func(m *mocks) {
				m.userRepo.EXPECT().
					GetByEmail(gomock.Any(), "test1").
					Return(nil, model.ErrUserNotFound)
			},
			args: args{
				email:    "test1",
				password: "password1",
			},
			wantRes: "",
			wantErr: model.ErrUserNotFound,
		},
		{
			name: "error.CreateToken",
			prepare: func(m *mocks) {
				m.userRepo.EXPECT().
					GetByEmail(gomock.Any(), "test1").
					Return(&model.User{
						UserID:   userID1,
						Email:    "test1",
						Password: makePasswordHash("password1"),
						UserRole: model.UserRoleEmployee,
					}, nil)
				m.tokenProvider.EXPECT().CreateToken("test1", userID1, model.UserRoleEmployee).Return("", assert.AnError)
			},
			args: args{
				email:    "test1",
				password: "password1",
			},
			wantRes: "",
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				userRepo:      NewMockuserRepo(ctrl),
				tokenProvider: NewMocktokenProvider(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.userRepo, m.tokenProvider)
			require.NoError(t, err)

			res, err := uc.Auth(context.Background(), tc.args.email, tc.args.password)
			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}

func makePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
