package reception_close

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

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMockreceptionClosing(ctrl), zap.NewNop())
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("error.first_nil", func(t *testing.T) {
		res, err := New(nil, zap.NewNop())
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("error.second_nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMockreceptionClosing(ctrl), nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestHandler_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockreceptionClosing(ctrl)

	pvzID, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a05")
	require.NoError(t, err)

	receptionID, err := model.ParseReceptionID("6451927e-846b-4c97-9924-cba818687a06")
	require.NoError(t, err)

	date := time.Date(2025, 4, 9, 20, 55, 59, 0, time.UTC)

	useCaseMock.EXPECT().
		CloseReception(gomock.Any(), pvzID).
		Return(model.Reception{
			ID:              receptionID,
			PVZID:           pvzID,
			ReceptionStatus: model.ReceptionStatusClose,
			ReceptedAt:      date,
		}, nil)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz/{pvzId}/close_last_reception", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	req.SetPathValue("pvzId", pvzID.UUID().String())
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"id": "6451927e-846b-4c97-9924-cba818687a06", "pvzId": "6451927e-846b-4c97-9924-cba818687a05", "dateTime": "2025-04-09T20:55:59Z", "status": "close"}`, w.Body.String())
}

func TestHandler_Handle_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockreceptionClosing(ctrl)

	pvzID, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a04")
	require.NoError(t, err)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz/{pvzId}/close_last_reception", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	req.SetPathValue("pvzId", pvzID.UUID().String())
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.JSONEq(t, `{"message": "only a user with the employee role can close reception"}`, w.Body.String())
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockreceptionClosing(ctrl)

	pvzID, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a03")
	require.NoError(t, err)

	useCaseMock.EXPECT().
		CloseReception(gomock.Any(), pvzID).
		Return(model.Reception{}, assert.AnError)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz/{pvzId}/close_last_reception", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	req.SetPathValue("pvzId", pvzID.UUID().String())
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, `{"message": "internal server error"}`, w.Body.String())
}

func TestHandler_Handle_InvalidPvzId(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockreceptionClosing(ctrl)

	pvzID1 := "6451927e-846b-4c97-9924-cba818687a0"

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz/{pvzId}/close_last_reception", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	req.SetPathValue("pvzId", pvzID1)
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "invalid pvzId"}`, w.Body.String())
}
