package registering

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMockuserRepo(ctrl))
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("error.first_nil", func(t *testing.T) {
		res, err := New(nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestUseCase_Register(t *testing.T) {
	type mocks struct {
		userRepo *MockuserRepo
	}
	type args struct {
		email    string
		password string
		role     model.UserRole
	}
	userID1 := model.NewUserID()
	passwordHash := makePasswordHash("password1")

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantRes *model.User
		wantErr error
	}{
		{
			name: "success.Register",
			prepare: func(m *mocks) {
				m.userRepo.EXPECT().
					Create(gomock.Any(), "test1", gomock.Any(), model.UserRoleEmployee).
					Return(&model.User{
						UserID:   userID1,
						Email:    "test1",
						Password: passwordHash,
						UserRole: model.UserRoleEmployee,
					}, nil)
			},
			args: args{
				email:    "test1",
				password: "password1",
				role:     model.UserRoleEmployee,
			},
			wantRes: &model.User{
				UserID:   userID1,
				Email:    "test1",
				Password: passwordHash,
				UserRole: model.UserRoleEmployee,
			},
			wantErr: nil,
		},
		{
			name: "error.CreateUser",
			prepare: func(m *mocks) {
				m.userRepo.EXPECT().
					Create(gomock.Any(), "test1", gomock.Any(), model.UserRoleEmployee).
					Return(nil, assert.AnError)
			},
			args: args{
				email:    "test1",
				password: "password1",
				role:     model.UserRoleEmployee,
			},
			wantRes: nil,
			wantErr: assert.AnError,
		},
		{
			name: "error.UserAlreadyExists",
			prepare: func(m *mocks) {
				m.userRepo.EXPECT().
					Create(gomock.Any(), "test1", gomock.Any(), model.UserRoleEmployee).
					Return(nil, model.ErrUserAlreadyExists)
			},
			args: args{
				email:    "test1",
				password: "password1",
				role:     model.UserRoleEmployee,
			},
			wantRes: nil,
			wantErr: model.ErrUserAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				userRepo: NewMockuserRepo(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.userRepo)
			require.NoError(t, err)

			res, err := uc.Register(context.Background(), tc.args.email, tc.args.password, tc.args.role)
			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}

func makePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
