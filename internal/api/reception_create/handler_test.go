package reception_create

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
	useCaseMock := NewMockreceptionCreating(ctrl)

	ID1, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a07")
	require.NoError(t, err)

	ID2, err := model.ParseReceptionID("6451927e-846b-4c97-9924-cba818687a06")
	require.NoError(t, err)

	date := time.Date(2025, 4, 9, 20, 55, 59, 0, time.UTC)

	useCaseMock.EXPECT().
		CreateReception(gomock.Any(), ID1).
		Return(model.Reception{
			ID:              ID2,
			PVZID:           ID1,
			ReceptionStatus: model.ReceptionStatusInProgress,
			ReceptedAt:      date,
		}, nil)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a07"}`)
	req := httptest.NewRequest(http.MethodPost, "/receptions/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	require.JSONEq(t, `{"id": "6451927e-846b-4c97-9924-cba818687a06", "pvzId": "6451927e-846b-4c97-9924-cba818687a07", "dateTime": "2025-04-09T20:55:59Z", "status": "in_progress"}`, w.Body.String())
}

func TestHandler_Handle_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockreceptionCreating(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a07"}`)
	req := httptest.NewRequest(http.MethodPost, "/receptions/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.JSONEq(t, `{"message": "only a user with the employee role can create reception"}`, w.Body.String())
}

func TestHandler_Handle_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockreceptionCreating(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a0"}`)
	req := httptest.NewRequest(http.MethodPost, "/receptions/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "could not bind request body"}`, w.Body.String())
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockreceptionCreating(ctrl)

	ID1, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a07")
	require.NoError(t, err)

	useCaseMock.EXPECT().
		CreateReception(gomock.Any(), ID1).
		Return(model.Reception{}, assert.AnError)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a07"}`)
	req := httptest.NewRequest(http.MethodPost, "/receptions/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	require.JSONEq(t, `{"message": "internal server error"}`, w.Body.String())
}
