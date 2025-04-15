package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	registry *prometheus.Registry

	httpRequestsTotal     *prometheus.CounterVec
	httpResponseTime      *prometheus.GaugeVec
	pvzCount              prometheus.Counter
	receptionCreatedCount prometheus.Counter
	productAddedCount     prometheus.Counter
}

func collectors(m *Metrics) []prometheus.Collector {
	return []prometheus.Collector{
		register(&m.httpRequestsTotal, prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество HTTP запросов",
		}, []string{"endpoint", "status_code"})),
		register(&m.httpResponseTime, prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "http_response_time",
			Help: "Время ответа на HTTP запросы в мс",
		}, []string{"endpoint", "status_code"})),
		register(&m.pvzCount, prometheus.NewCounter(prometheus.CounterOpts{
			Name: "pvz_registered_count",
			Help: "Количество созданных ПВЗ",
		})),
		register(&m.receptionCreatedCount, prometheus.NewCounter(prometheus.CounterOpts{
			Name: "reception_created_count",
			Help: "Количество созданных приёмок заказов",
		})),
		register(&m.productAddedCount, prometheus.NewCounter(prometheus.CounterOpts{
			Name: "product_added_count",
			Help: "Количество добавленных товаров",
		})),
	}
}

func New() (*Metrics, error) {
	registry := prometheus.NewRegistry()
	m := &Metrics{registry: registry}

	for _, collector := range collectors(m) {
		err := m.registry.Register(collector)
		if err != nil {
			return nil, fmt.Errorf("registry.Register: %w", err)
		}
	}

	return m, nil
}

func register[T prometheus.Collector](field *T, collector T) T {
	*field = collector
	return collector
}

func (m *Metrics) PVZRegisteredCountInc() {
	m.pvzCount.Inc()
}

func (m *Metrics) ReceptionCreatedCountInc() {
	m.receptionCreatedCount.Inc()
}

func (m *Metrics) ProductAddedCountInc() {
	m.productAddedCount.Inc()
}
