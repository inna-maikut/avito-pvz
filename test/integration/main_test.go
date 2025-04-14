//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-pvz/internal/api"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/config"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

func setUp() {
	_ = config.Load()
}

func dummyLogin(t *testing.T, role model.UserRole) string {
	resp := apiPost(t, "/dummyLogin", "", api.PostDummyLoginJSONBody{
		Role: api.PostDummyLoginJSONBodyRole(role),
	})
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, http.StatusOK, resp.StatusCode, "body: "+string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	require.NoError(t, err)

	return strings.Trim(string(bodyBytes), "\"\n")
}

// apiGet path should start with slash
func apiGet(t *testing.T, path, token string) *http.Response {
	t.Helper()

	url := "http://localhost:" + os.Getenv("SERVER_PORT") + path
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// apiPost path should start with slash
func apiPost[In any](t *testing.T, path, token string, in In) *http.Response {
	t.Helper()

	inStr, err := json.Marshal(in)
	require.NoError(t, err)

	url := "http://localhost:" + os.Getenv("SERVER_PORT") + path
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(inStr))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func parseJSON[Out any](t *testing.T, resp *http.Response) Out {
	var out Out

	bodyBytes, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	require.NoError(t, err)

	err = json.Unmarshal(bodyBytes, &out)
	require.NoError(t, err)

	return out
}

func assertStatus(t *testing.T, resp *http.Response, expected int) {
	if resp.StatusCode == expected {
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode, "body: "+string(bodyBytes))
}
