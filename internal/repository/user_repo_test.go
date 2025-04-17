//go:build integration

package repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestNewUserRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := NewUserRepository(&sqlx.DB{}, &trmsqlx.CtxGetter{})
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("error.first_nil", func(t *testing.T) {
		res, err := NewUserRepository(nil, &trmsqlx.CtxGetter{})
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.second_nil", func(t *testing.T) {
		res, err := NewUserRepository(&sqlx.DB{}, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func Test_GetByEmail(t *testing.T) {
	db := setUp(t)
	repo, err := NewUserRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	userID := model.NewUserID()

	type args struct {
		email string
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		wantRes *model.User
		wantErr error
	}{
		{
			name: "found",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM users where email = $1`, "get-by-email-1")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO users (id, email, password, user_role)
					VALUES ($1, $2, $3, $4)`, userID, "get-by-email-1", "password", model.UserRoleEmployee)
				require.NoError(t, err)
			},
			args: args{
				email: "get-by-email-1",
			},
			wantRes: &model.User{
				UserID:   userID,
				Email:    "get-by-email-1",
				Password: "password",
				UserRole: model.UserRoleEmployee,
			},
			wantErr: nil,
		},
		{
			name: "not_found",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM users where email = $1`, "get-by-email-2")
				require.NoError(t, err)
			},
			args: args{
				email: "get-by-email-2",
			},
			wantRes: nil,
			wantErr: model.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.GetByEmail(context.Background(), tc.args.email)

			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.wantRes, res)
		})
	}
}
func Test_Create(t *testing.T) {
	db := setUp(t)
	repo, err := NewUserRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	userID := model.NewUserID()

	type args struct {
		email        string
		passwordHash string
		role         model.UserRole
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		wantRes *model.User
		wantErr error
	}{
		{
			name: "create",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM users where email = $1`, "get-by-email-2")
				require.NoError(t, err)
			},
			args: args{
				email:        "get-by-email-2",
				passwordHash: "password",
				role:         model.UserRoleModerator,
			},
			wantRes: &model.User{
				UserID:   userID,
				Email:    "get-by-email-2",
				Password: "password",
				UserRole: model.UserRoleModerator,
			},
			wantErr: nil,
		},
		{
			name: "error.UserAlreadyExists",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM users where email = $1`, "get-by-email-2")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO users (id, email, password, user_role)
					VALUES ($1, $2, $3, $4)`, userID, "get-by-email-2", "password", model.UserRoleModerator)
				require.NoError(t, err)
			},
			args: args{
				email:        "get-by-email-2",
				passwordHash: "password",
				role:         model.UserRoleModerator,
			},
			wantRes: nil,
			wantErr: model.ErrUserAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.Create(context.Background(), tc.args.email, tc.args.passwordHash, tc.args.role)

			require.ErrorIs(t, err, tc.wantErr)
			if res != nil {
				res.UserID = userID // can't validate id, because it's generated in DataBase
			}
			require.Equal(t, tc.wantRes, res)
		})
	}
}
