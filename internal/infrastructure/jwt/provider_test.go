package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestProvider_CreateToken(t *testing.T) {
	var secret = "secret"

	type args struct {
		email  string
		userID int64
		role   model.UserRole
	}
	tests := []struct {
		name    string
		args    args
		check   func(t *testing.T, got string)
		wantErr bool
	}{
		{
			name: "success.moderator",
			args: args{
				email:  "test@test.com",
				userID: 1,
				role:   model.UserRoleModerator,
			},
			check: func(t *testing.T, got string) {
				require.NotEmpty(t, got)
				token, err := jwt.Parse(got, func(_ *jwt.Token) (any, error) {
					return []byte(secret), nil
				})
				require.NoError(t, err)

				claims, _ := token.Claims.(jwt.MapClaims)
				assert.InDelta(t, claims["exp"], time.Now().Add(time.Hour*72).Unix(), 10)
				delete(claims, "exp")
				assert.Equal(t, jwt.MapClaims{
					"email":  "test@test.com",
					"userID": float64(1),
					"role":   "moderator",
				}, claims)
			},
			wantErr: false,
		},
		{
			name: "success.employee",
			args: args{
				email:  "test@test.com",
				userID: 2,
				role:   model.UserRoleEmployee,
			},
			check: func(t *testing.T, got string) {
				require.NotEmpty(t, got)
				token, err := jwt.Parse(got, func(_ *jwt.Token) (any, error) {
					return []byte(secret), nil
				})
				require.NoError(t, err)

				claims, _ := token.Claims.(jwt.MapClaims)
				assert.InDelta(t, claims["exp"], time.Now().Add(time.Hour*72).Unix(), 10)
				delete(claims, "exp")
				assert.Equal(t, jwt.MapClaims{
					"email":  "test@test.com",
					"userID": float64(2),
					"role":   "employee",
				}, claims)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(secret)
			got, err := p.CreateToken(tt.args.email, tt.args.userID, tt.args.role)
			require.Equal(t, err != nil, tt.wantErr)
			tt.check(t, got)
		})
	}
}

func TestProvider_ParseToken(t *testing.T) {
	var secret = "secret"
	exp := time.Now().Add(time.Hour * 72).Unix()

	tests := []struct {
		name     string
		getToken func(t *testing.T) string
		want     model.TokenInfo
		wantErr  bool
	}{
		{
			name: "success.moderator",
			getToken: func(t *testing.T) string {

				claims := jwt.MapClaims{
					"email":  "email",
					"userID": 123,
					"role":   "moderator",
					"exp":    exp,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, err := token.SignedString([]byte(secret))
				require.NoError(t, err)

				return tokenStr
			},
			want: model.TokenInfo{
				Email:    "email",
				UserID:   123,
				UserRole: model.UserRoleModerator,
			},
			wantErr: false,
		},
		{
			name: "success.employee",
			getToken: func(t *testing.T) string {

				claims := jwt.MapClaims{
					"email":  "email",
					"userID": 123,
					"role":   "employee",
					"exp":    exp,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, err := token.SignedString([]byte(secret))
				require.NoError(t, err)

				return tokenStr
			},
			want: model.TokenInfo{
				Email:    "email",
				UserID:   123,
				UserRole: model.UserRoleEmployee,
			},
			wantErr: false,
		},
		{
			name: "err.wrong_sign",
			getToken: func(t *testing.T) string {

				claims := jwt.MapClaims{
					"email":  "email",
					"userID": 123,
					"role":   "employee",
					"exp":    exp,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, err := token.SignedString([]byte(secret + "1"))
				require.NoError(t, err)

				return tokenStr
			},
			want:    model.TokenInfo{},
			wantErr: true,
		},
		{
			name: "err.wrong_exp",
			getToken: func(t *testing.T) string {
				claims := jwt.MapClaims{
					"email":  "email",
					"userID": 123,
					"role":   "employee",
					"exp":    time.Now().Add(-time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, err := token.SignedString([]byte(secret))
				require.NoError(t, err)

				return tokenStr
			},
			want:    model.TokenInfo{},
			wantErr: true,
		},
		{
			name: "err.wrong_claims.user_id",
			getToken: func(t *testing.T) string {
				claims := jwt.MapClaims{
					"email":  "email",
					"userID": "123",
					"role":   "employee",
					"exp":    exp,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, err := token.SignedString([]byte(secret))
				require.NoError(t, err)

				return tokenStr
			},
			want:    model.TokenInfo{},
			wantErr: true,
		},
		{
			name: "err.wrong_claims.email",
			getToken: func(t *testing.T) string {
				claims := jwt.MapClaims{
					"email":  123,
					"userID": 123,
					"role":   "employee",
					"exp":    exp,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, err := token.SignedString([]byte(secret))
				require.NoError(t, err)

				return tokenStr
			},
			want:    model.TokenInfo{},
			wantErr: true,
		},
		{
			name: "err.wrong_claims.no_role",
			getToken: func(t *testing.T) string {
				claims := jwt.MapClaims{
					"email":  "email",
					"userID": 123,
					"exp":    exp,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, err := token.SignedString([]byte(secret))
				require.NoError(t, err)

				return tokenStr
			},
			want:    model.TokenInfo{},
			wantErr: true,
		},
		{
			name: "err.wrong_claims.invalid_role",
			getToken: func(t *testing.T) string {
				claims := jwt.MapClaims{
					"email":  "email",
					"userID": 123,
					"role":   "invalid",
					"exp":    exp,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, err := token.SignedString([]byte(secret))
				require.NoError(t, err)

				return tokenStr
			},
			want:    model.TokenInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(secret)
			token := tt.getToken(t)
			got, err := p.ParseToken(token)
			require.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
