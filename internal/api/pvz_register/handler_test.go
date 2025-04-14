package pvz_register

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestHandler_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzRegistering(ctrl)

	ID1, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a07")
	require.NoError(t, err)

	date := time.Date(2025, 4, 9, 20, 55, 59, 0, time.UTC)

	useCaseMock.EXPECT().
		RegisterPVZ(gomock.Any(), "Москва").
		Return(model.PVZ{
			ID:           ID1,
			City:         "Москва",
			RegisteredAt: date,
		}, nil)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"city": "Москва"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/pvz", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	require.JSONEq(t, `{"city": "Москва", "id": "6451927e-846b-4c97-9924-cba818687a07", "registrationDate": "2025-04-09T20:55:59Z"}`, w.Body.String())
}

func TestHandler_Handle_InvalidCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzRegistering(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"city": "test3"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/pvz", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "invalid city"}`, w.Body.String())
}

func TestHandler_Handle_NotModerator(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzRegistering(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"city": "Москва"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/pvz", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.JSONEq(t, `{"message": "only a user with the moderator role can create a pickup point in the system"}`, w.Body.String())
}

func TestHandler_Handle_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzRegistering(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"city": 0}`)
	req := httptest.NewRequest(http.MethodPost, "/api/pvz", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "could not bind request body"}`, w.Body.String())
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzRegistering(ctrl)

	useCaseMock.EXPECT().
		RegisterPVZ(gomock.Any(), "Москва").
		Return(model.PVZ{}, assert.AnError)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"city": "Москва"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/pvz", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	require.JSONEq(t, `{"message": "internal server error"}`, w.Body.String())
}
