package jwt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestContext(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := ContextWithTokenInfo(context.Background(), model.TokenInfo{UserID: 1})

		tokenInfo := TokenInfoFromContext(ctx)
		assert.Equal(t, model.TokenInfo{UserID: 1}, tokenInfo)
	})

	t.Run("empty", func(t *testing.T) {
		tokenInfo := TokenInfoFromContext(context.Background())
		assert.Equal(t, model.TokenInfo{}, tokenInfo)
	})

	t.Run("empty_if_has_key_empty_struct", func(t *testing.T) {
		tokenInfo := TokenInfoFromContext(context.WithValue(context.Background(), struct{}{}, model.TokenInfo{UserID: 1}))
		assert.Equal(t, model.TokenInfo{}, tokenInfo)
	})
}

func TestGetJWSFromRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "asdf")

		res, err := GetJWSFromRequest(r)
		require.NoError(t, err)
		require.Equal(t, "asdf", res)
	})

	t.Run("error.no_header", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		_, err := GetJWSFromRequest(r)
		require.ErrorIs(t, err, ErrNoAuthHeader)
	})

	t.Run("error.empty_header", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "")

		_, err := GetJWSFromRequest(r)
		require.ErrorIs(t, err, ErrNoAuthHeader)
	})
}

func TestAuthenticate(t *testing.T) {
	type mocks struct {
		tokenProvider *MocktokenProvider
	}
	type args struct {
		ctx   context.Context
		input *openapi3filter.AuthenticationInput
	}

	tests := []struct {
		name    string
		args    func(t *testing.T) args
		prepare func(t *testing.T, m *mocks)
		check   func(t *testing.T, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: func(t *testing.T) args {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "asdf")
				return args{
					ctx: context.Background(),
					input: &openapi3filter.AuthenticationInput{
						RequestValidationInput: &openapi3filter.RequestValidationInput{
							Request: req,
						},
						SecuritySchemeName: "bearerAuth",
					},
				}
			},
			prepare: func(t *testing.T, m *mocks) {
				m.tokenProvider.EXPECT().ParseToken("asdf").Return(model.TokenInfo{UserID: 1}, nil)
			},
			check: func(t *testing.T, args args) {
				tokenInfo := TokenInfoFromContext(args.input.RequestValidationInput.Request.Context())
				assert.Equal(t, model.TokenInfo{UserID: 1}, tokenInfo)
			},
			wantErr: false,
		},
		{
			name: "error.invalid_security_scheme",
			args: func(t *testing.T) args {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "asdf")
				return args{
					input: &openapi3filter.AuthenticationInput{
						RequestValidationInput: &openapi3filter.RequestValidationInput{
							Request: req,
						},
						SecuritySchemeName: "bearerAuth1",
					},
				}
			},
			prepare: func(_ *testing.T, _ *mocks) {},
			check:   func(_ *testing.T, _ args) {},
			wantErr: true,
		},
		{
			name: "error.parse_token",
			args: func(t *testing.T) args {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "asdf")
				return args{
					input: &openapi3filter.AuthenticationInput{
						RequestValidationInput: &openapi3filter.RequestValidationInput{
							Request: req,
						},
						SecuritySchemeName: "bearerAuth",
					},
				}
			},
			prepare: func(t *testing.T, m *mocks) {
				m.tokenProvider.EXPECT().ParseToken("asdf").Return(model.TokenInfo{}, assert.AnError)
			},
			check:   func(_ *testing.T, _ args) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				tokenProvider: NewMocktokenProvider(ctrl),
			}

			tt.prepare(t, m)

			a := tt.args(t)
			err := Authenticate(a.ctx, m.tokenProvider, a.input)
			require.Equal(t, err != nil, tt.wantErr)
			tt.check(t, a)
		})
	}
}
