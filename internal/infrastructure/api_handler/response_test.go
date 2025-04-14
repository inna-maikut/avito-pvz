package api_handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	InternalError(w, "my description")

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, `{"message": "my description"}`, w.Body.String())
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequest(w, "my description")

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"message": "my description"}`, w.Body.String())
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	Unauthorized(w, "my description")

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.JSONEq(t, `{"message": "my description"}`, w.Body.String())
}

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	Forbidden(w, "my description")

	require.Equal(t, http.StatusForbidden, w.Code)
	require.JSONEq(t, `{"message": "my description"}`, w.Body.String())
}

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	OK(w, "my description")

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `"my description"`, w.Body.String())
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	Created(w, "my description")

	require.Equal(t, http.StatusCreated, w.Code)
	require.JSONEq(t, `"my description"`, w.Body.String())
}
