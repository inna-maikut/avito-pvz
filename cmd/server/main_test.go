//go:build integration

package main

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal/infrastructure/config"
)

func TestMetrics_RunHTTPServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Config{
		ServerPort: 9003,
	}
	done := make(chan struct{})
	go func() {
		mux := http.NewServeMux()
		mux.Handle("POST /dummyLogin", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		runHTTPServer(ctx, mux, cfg, zap.NewNop())
		done <- struct{}{}
	}()

	time.Sleep(time.Millisecond)

	dummyLoginBody := `{"role": "moderator"}`
	require.Eventually(t, func() bool {
		_, err := http.Post("http://localhost:9003/dummyLogin", "application/json", strings.NewReader(dummyLoginBody))

		return err == nil
	}, 100*time.Millisecond, 100*time.Nanosecond)

	resp, err := http.Post("http://localhost:9003/dummyLogin", "application/json", strings.NewReader(dummyLoginBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = resp.Body.Close()
	require.NoError(t, err)
	respString := string(respBody)

	assert.Contains(t, respString, "ok")

	cancel()

	require.Eventually(t, func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, 100*time.Millisecond, 100*time.Nanosecond)
}
