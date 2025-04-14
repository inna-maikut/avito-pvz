package product_add

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
		res, err := New(NewMockproductAdding(ctrl), zap.NewNop())
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
		res, err := New(NewMockproductAdding(ctrl), nil)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestHandler_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductAdding(ctrl)

	pvzID, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a05")
	require.NoError(t, err)
	receptionID, err := model.ParseReceptionID("6451927e-846b-4c97-9924-cba818687a04")
	require.NoError(t, err)
	productID, err := model.ParseProductID("6451927e-846b-4c97-9924-cba818687a03")
	require.NoError(t, err)

	date := time.Date(2025, 4, 9, 20, 55, 59, 0, time.UTC)

	useCaseMock.EXPECT().
		AddProduct(gomock.Any(), pvzID, model.ProductCategoryElectronics).
		Return(model.Product{
			ID:          productID,
			ReceptionID: receptionID,
			Category:    model.ProductCategoryElectronics,
			AddedAt:     date,
		}, nil)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a05", "type": "электроника"}`)
	req := httptest.NewRequest(http.MethodPost, "/products/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	require.JSONEq(t, `{"id": "6451927e-846b-4c97-9924-cba818687a03", "receptionId": "6451927e-846b-4c97-9924-cba818687a04", "dateTime": "2025-04-09T20:55:59Z", "type": "электроника"}`, w.Body.String())
}

func TestHandler_Handle_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductAdding(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a05", "type": "электроника"}`)
	req := httptest.NewRequest(http.MethodPost, "/products/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.JSONEq(t, `{"message": "only a user with the employee role can add product"}`, w.Body.String())
}

func TestHandler_Handle_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductAdding(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a0", "type": "электроника"}`)
	req := httptest.NewRequest(http.MethodPost, "/products/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "could not bind request body"}`, w.Body.String())
}

func TestHandler_Handle_InvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductAdding(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a05", "type": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/products/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "invalid type"}`, w.Body.String())
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductAdding(ctrl)

	useCaseMock.EXPECT().
		AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.Product{}, assert.AnError)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"pvzId": "6451927e-846b-4c97-9924-cba818687a05", "type": "электроника"}`)
	req := httptest.NewRequest(http.MethodPost, "/products/", bytes.NewReader(validData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	require.JSONEq(t, `{"message": "internal server error"}`, w.Body.String())
}
