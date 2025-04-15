package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal/infrastructure/config"
)

func (m *Metrics) RunHTTPServer(ctx context.Context, cfg config.Config, logger *zap.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registry}))

	s := http.Server{
		Handler:           mux,
		Addr:              "0.0.0.0:" + strconv.Itoa(cfg.MetricsServerPort),
		ReadHeaderTimeout: time.Second,
	}
	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := s.Shutdown(shutdownCtx); err != nil {
			err = fmt.Errorf("HTTP metrics server shutdown error: %w", err)
			logger.Error("HTTP metrics server shutdown error", zap.Error(err))
		}
	}()

	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("HTTP metrics server shutdown error", zap.Error(err))
	}
}
