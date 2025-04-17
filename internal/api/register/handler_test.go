package register

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
		res, err := New(NewMockregistering(ctrl), zap.NewNop())
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
		res, err := New(NewMockregistering(ctrl), nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestHandler_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	registeringMock := NewMockregistering(ctrl)

	userID := model.NewUserID()

	registeringMock.EXPECT().
		Register(gomock.Any(), "email1@gmail.com", "password1", model.UserRoleModerator).
		Return(&model.User{
			UserID:   userID,
			Email:    "email1@gmail.com",
			Password: "password1",
			UserRole: model.UserRoleModerator,
		}, nil)

	handler, err := New(registeringMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"email": "email1@gmail.com", "password" : "password1", "role":"moderator"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.JSONEq(t, `{"id": "`+userID.UUID().String()+`", "email": "email1@gmail.com", "role":"moderator"}`, w.Body.String())
}

func TestHandler_Handle_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	registeringMock := NewMockregistering(ctrl)

	registeringMock.EXPECT().
		Register(gomock.Any(), "email1@gmail.com", "password1", model.UserRoleModerator).
		Return(nil, model.ErrUserAlreadyExists)

	handler, err := New(registeringMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"email": "email1@gmail.com", "password" : "password1", "role":"moderator"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "user already exists"}`, w.Body.String())
}
func TestHandler_Handle_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	registeringMock := NewMockregistering(ctrl)

	handler, err := New(registeringMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"email": "email1@gmail.com", "password" : "password1", "role":""}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "invalid role"}`, w.Body.String())
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	registeringMock := NewMockregistering(ctrl)

	registeringMock.EXPECT().
		Register(gomock.Any(), "email1@gmail.com", "password1", model.UserRoleModerator).
		Return(nil, assert.AnError)

	handler, err := New(registeringMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"email": "email1@gmail.com", "password" : "password1", "role":"moderator"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, `{"message": "internal server error"}`, w.Body.String())
}

func TestHandler_Handle_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	registeringMock := NewMockregistering(ctrl)

	handler, err := New(registeringMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"email": "email1@gmail", "password" : "password1", "role":"moderator"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "could not bind request body"}`, w.Body.String())
}
