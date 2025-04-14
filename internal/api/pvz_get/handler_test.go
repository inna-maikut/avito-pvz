package pvz_get

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
	useCaseMock := NewMockpvzListGetting(ctrl)

	pvzID1, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a11")
	require.NoError(t, err)
	pvzID2, err := model.ParsePVZID("6451927e-846b-4c97-9924-cba818687a12")
	require.NoError(t, err)
	receptionID1, err := model.ParseReceptionID("6451927e-846b-4c97-9924-cba818687a21")
	require.NoError(t, err)
	receptionID2, err := model.ParseReceptionID("6451927e-846b-4c97-9924-cba818687a22")
	require.NoError(t, err)
	productID1, err := model.ParseProductID("6451927e-846b-4c97-9924-cba818687a31")
	require.NoError(t, err)
	productID2, err := model.ParseProductID("6451927e-846b-4c97-9924-cba818687a32")
	require.NoError(t, err)
	productID3, err := model.ParseProductID("6451927e-846b-4c97-9924-cba818687a33")
	require.NoError(t, err)

	from := time.Date(2025, 4, 9, 20, 55, 59, 0, time.UTC)
	to := from.Add(1 * time.Hour)

	useCaseMock.EXPECT().
		GetPVZList(gomock.Any(), &from, &to, int64(1), int64(30)).
		Return(model.PVZList{
			PVZs: []model.PVZ{
				{
					ID:           pvzID1,
					City:         "Москва",
					RegisteredAt: from,
				},
				{
					ID:           pvzID2,
					City:         "Санкт-Петербург",
					RegisteredAt: from.Add(10 * time.Second),
				},
			},
			Receptions: []model.Reception{
				{
					ID:              receptionID1,
					PVZID:           pvzID1,
					ReceptionStatus: model.ReceptionStatusClose,
					ReceptedAt:      from.Add(1 * time.Second),
				},
				{
					ID:              receptionID2,
					PVZID:           pvzID2,
					ReceptionStatus: model.ReceptionStatusInProgress,
					ReceptedAt:      from.Add(12 * time.Second),
				},
			},
			Products: []model.Product{
				{
					ID:          productID1,
					ReceptionID: receptionID1,
					Category:    model.ProductCategoryElectronics,
					AddedAt:     from.Add(1 * time.Minute),
				},
				{
					ID:          productID2,
					ReceptionID: receptionID2,
					Category:    model.ProductCategoryClothes,
					AddedAt:     from.Add(2 * time.Minute),
				},
				{
					ID:          productID3,
					ReceptionID: receptionID2,
					Category:    model.ProductCategoryShoes,
					AddedAt:     from.Add(3 * time.Minute),
				},
			},
		}, nil)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/pvz?startDate=2025-04-09T20:55:59.000Z&endDate=2025-04-09T21:55:59.000Z&page=1&limit=30", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `[{
		"pvz": {
			"id": "6451927e-846b-4c97-9924-cba818687a11",
			"city": "Москва",
			"registrationDate": "2025-04-09T20:55:59Z"
		},
		"receptions": [
			{
				"reception": {
					"id": "6451927e-846b-4c97-9924-cba818687a21",
					"pvzId": "6451927e-846b-4c97-9924-cba818687a11",
					"status": "close",
					"dateTime": "2025-04-09T20:56:00Z"
				},
				"products": [
					{
						"id": "6451927e-846b-4c97-9924-cba818687a31",
						"receptionId": "6451927e-846b-4c97-9924-cba818687a21",
						"type": "электроника",
						"dateTime": "2025-04-09T20:56:59Z"
					}
				]
			}
		]
	}, {
		"pvz": {
			"id": "6451927e-846b-4c97-9924-cba818687a12",
			"city": "Санкт-Петербург",
			"registrationDate": "2025-04-09T20:56:09Z"
		},
		"receptions": [
			{
				"reception": {
					"id": "6451927e-846b-4c97-9924-cba818687a22",
					"pvzId": "6451927e-846b-4c97-9924-cba818687a12",
					"status": "in_progress",
					"dateTime": "2025-04-09T20:56:11Z"
				},
				"products": [
					{
						"id": "6451927e-846b-4c97-9924-cba818687a32",
						"receptionId": "6451927e-846b-4c97-9924-cba818687a22",
						"type": "одежда",
						"dateTime": "2025-04-09T20:57:59Z"
					},
					{
						"id": "6451927e-846b-4c97-9924-cba818687a33",
						"receptionId": "6451927e-846b-4c97-9924-cba818687a22",
						"type": "обувь",
						"dateTime": "2025-04-09T20:58:59Z"
					}
				]
			}
		]
	}]`, w.Body.String())
}

func TestHandler_Handle_InvalidStartDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzListGetting(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/pvz?startDate=2025-04-091", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "validation query: parse start date: parsing time \"2025-04-091\": extra text: \"1\""}`, w.Body.String())
}

func TestHandler_Handle_InvalidEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzListGetting(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/pvz?endDate=2025-04-091", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "validation query: parse end date: parsing time \"2025-04-091\": extra text: \"1\""}`, w.Body.String())
}

func TestHandler_Handle_InvalidLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzListGetting(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/pvz?limit=31", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "validation query: limit must be not greater than 30"}`, w.Body.String())
}

func TestHandler_Handle_InvalidPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzListGetting(ctrl)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/pvz?page=0", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "validation query: page must be greater than zero"}`, w.Body.String())
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCaseMock := NewMockpvzListGetting(ctrl)

	useCaseMock.EXPECT().
		GetPVZList(gomock.Any(), nil, nil, int64(1), int64(10)).
		Return(model.PVZList{}, assert.AnError)

	handler, err := New(useCaseMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/pvz", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		UserRole: model.UserRoleModerator,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	require.JSONEq(t, `{"message": "internal server error"}`, w.Body.String())
}
