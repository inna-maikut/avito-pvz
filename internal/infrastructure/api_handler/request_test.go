package api_handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type Val struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name       string
		body       string
		wantOk     bool
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			body:       `{"name":"test"}`,
			wantOk:     true,
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
		{
			name:       "error.parse_json",
			body:       `{"name":1}`,
			wantOk:     false,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"could not bind request body"}`,
		},
		{
			name:       "error.empty_body",
			body:       "",
			wantOk:     false,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"message":"could not bind request body"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Val

			r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			gotOK := Parse(r, w, &v)
			assert.Equal(t, tt.wantOk, gotOK)
			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			}
		})
	}
}
