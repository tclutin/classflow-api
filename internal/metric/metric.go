package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var httpRequestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "classflow_api",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total HTTP requests",
	},
	[]string{"endpoint"},
)

func init() {
	prometheus.MustRegister(httpRequestsCounter)
}

func IncRequestCounter(endpoint string) {
	httpRequestsCounter.WithLabelValues(endpoint).Inc()
}
