//go:build integration

package metrics

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

	m, err := New()
	require.NoError(t, err)
	cfg := config.Config{
		MetricsServerPort: 9001,
	}
	done := make(chan struct{})
	go func() {
		m.RunHTTPServer(ctx, cfg, zap.NewNop())
		done <- struct{}{}
	}()
	done2 := make(chan struct{})
	go func() {
		mux := http.NewServeMux()
		mux.Handle("POST /dummyLogin", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		s := http.Server{
			Handler:           m.HTTPServerMW(mux),
			Addr:              "0.0.0.0:9002",
			ReadHeaderTimeout: time.Second,
		}
		go func() {
			<-ctx.Done()
			_ = s.Shutdown(ctx)
		}()
		_ = s.ListenAndServe()
		done2 <- struct{}{}
	}()

	time.Sleep(time.Millisecond)

	require.Eventually(t, func() bool {
		_, err2 := http.Get("http://localhost:9001/metrics")

		return err2 == nil
	}, 100*time.Millisecond, 100*time.Nanosecond)

	dummyLoginBody := `{"role": "moderator"}`
	require.Eventually(t, func() bool {
		_, err2 := http.Post("http://localhost:9002/dummyLogin", "application/json", strings.NewReader(dummyLoginBody))

		return err2 == nil
	}, 100*time.Millisecond, 100*time.Nanosecond)

	resp, err := http.Post("http://localhost:9002/dummyLogin", "application/json", strings.NewReader(dummyLoginBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	m.PVZRegisteredCountInc()
	m.ReceptionCreatedCountInc()
	m.ReceptionCreatedCountInc()
	m.ProductAddedCountInc()
	m.ProductAddedCountInc()
	m.ProductAddedCountInc()

	resp, err = http.Get("http://localhost:9001/metrics")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = resp.Body.Close()
	require.NoError(t, err)
	respString := string(respBody)

	assert.Contains(t, respString, "pvz_registered_count 1")
	assert.Contains(t, respString, "reception_created_count 2")
	assert.Contains(t, respString, "product_added_count 3")
	assert.Contains(t, respString, "http_requests_total{endpoint=\"POST__/dummyLogin\",status_code=\"200\"} 2")
	assert.Contains(t, respString, "http_response_time{endpoint=\"POST__/dummyLogin\",status_code=\"200\"} ")

	cancel()

	require.Eventually(t, func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, 100*time.Millisecond, 100*time.Nanosecond)
	require.Eventually(t, func() bool {
		select {
		case <-done2:
			return true
		default:
			return false
		}
	}, 100*time.Millisecond, 100*time.Nanosecond)
}
