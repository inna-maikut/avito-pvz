package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestCreateNoAuthMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		check   func(t *testing.T, mw func(next http.Handler) http.Handler)
		wantErr bool
	}{
		{
			name: "validation_required",
			check: func(t *testing.T, mw func(next http.Handler) http.Handler) {
				called := false
				next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
					called = true
				})
				handler := mw(next)

				r := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(nil))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)

				require.False(t, called)
				require.Equal(t, "request body has an error: value is required but missing\n", w.Body.String())
			},
		},
		{
			name: "pass_no_auth",
			check: func(t *testing.T, mw func(next http.Handler) http.Handler) {
				called := false
				next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
					called = true
				})
				handler := mw(next)

				r := httptest.NewRequest(http.MethodGet, "/pvz", bytes.NewReader(nil))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)

				require.True(t, called)
				require.Equal(t, "", w.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateNoAuthMiddleware()
			require.NoError(t, err)

			tt.check(t, got)
		})
	}
}

func TestCreateAuthMiddleware(t *testing.T) {
	type mocks struct {
		tokenProvider *MocktokenProvider
	}

	tests := []struct {
		name    string
		prepare func(t *testing.T, m *mocks)
		check   func(t *testing.T, mw func(next http.Handler) http.Handler)
		wantErr bool
	}{
		{
			name: "validation_required",
			prepare: func(_ *testing.T, m *mocks) {
				m.tokenProvider.EXPECT().ParseToken("asdf").Return(model.TokenInfo{}, nil)
			},
			check: func(t *testing.T, mw func(next http.Handler) http.Handler) {
				called := false
				next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
					called = true
				})
				handler := mw(next)

				r := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(nil))
				r.Header.Set("Authorization", "asdf")
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)

				assert.False(t, called)
				assert.Equal(t, "request body has an error: value is required but missing\n", w.Body.String())
			},
		},
		{
			name:    "forbidden_no_header",
			prepare: func(_ *testing.T, _ *mocks) {},
			check: func(t *testing.T, mw func(next http.Handler) http.Handler) {
				called := false
				next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
					called = true
				})
				handler := mw(next)

				r := httptest.NewRequest(http.MethodGet, "/pvz", bytes.NewReader(nil))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)

				assert.False(t, called)
				assert.Equal(t, "security requirements failed: getting jws: authorization header is missing\n", w.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				tokenProvider: NewMocktokenProvider(ctrl),
			}

			tt.prepare(t, m)

			got, err := CreateAuthMiddleware(m.tokenProvider)
			require.NoError(t, err)

			tt.check(t, got)
		})
	}
}
