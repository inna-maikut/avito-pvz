package product_remove_last

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestHandler_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductRemoving(ctrl)

	pvzID, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a05")
	require.NoError(t, err)

	useCaseMock.EXPECT().
		RemoveLastProduct(gomock.Any(), pvzID).
		Return(nil)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz/{pvzId}/delete_last_product", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	req.SetPathValue("pvzId", pvzID.UUID().String())
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Handle_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductRemoving(ctrl)

	pvzID, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a04")
	require.NoError(t, err)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz/{pvzId}/delete_last_product", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	req.SetPathValue("pvzId", pvzID.UUID().String())
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductRemoving(ctrl)

	pvzID, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a03")
	require.NoError(t, err)

	useCaseMock.EXPECT().
		RemoveLastProduct(gomock.Any(), pvzID).
		Return(assert.AnError)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz/{pvzId}/delete_last_product", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	req.SetPathValue("pvzId", pvzID.UUID().String())
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_Handle_InvalidPvzId(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockproductRemoving(ctrl)

	pvzID1 := "6451927e-846b-4c97-9924-cba818687a0"

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz/{pvzId}/delete_last_product", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleEmployee,
	}))
	req.SetPathValue("pvzId", pvzID1)
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
