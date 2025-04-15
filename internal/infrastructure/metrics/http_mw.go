package metrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (m *Metrics) HTTPServerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		loggingW := newLoggingResponseWriter(w)

		next.ServeHTTP(loggingW, r)

		labelValues := []string{httpEndpoint(r), strconv.Itoa(loggingW.statusCode)}
		m.httpRequestsTotal.WithLabelValues(labelValues...).Inc()
		m.httpResponseTime.WithLabelValues(labelValues...).Set(float64(time.Since(startTime)) / float64(time.Millisecond))
	})
}

func httpEndpoint(r *http.Request) string {
	return strings.ReplaceAll(r.Pattern, " ", "__")
}
