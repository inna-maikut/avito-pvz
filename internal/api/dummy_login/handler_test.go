package dummy_login

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		res, err := New(NewMockauthenticating(ctrl), zap.NewNop())
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
		res, err := New(NewMockauthenticating(ctrl), nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestHandler_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	authenticatingMock.EXPECT().
		Auth(gomock.Any(), model.UserRoleModerator).
		Return("token1", nil)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"role": "moderator"}`)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `"token1"`, w.Body.String())
}

func TestHandler_Handle_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"role": "moderator1"}`)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "invalid role"}`, w.Body.String())
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	authenticatingMock.EXPECT().
		Auth(gomock.Any(), model.UserRoleModerator).
		Return("", assert.AnError)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"role": "moderator"}`)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, `{"message": "internal server error"}`, w.Body.String())
}

func TestHandler_Handle_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"role": 1}`)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "could not bind request body"}`, w.Body.String())
}
