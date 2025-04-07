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

var scheduleRequestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "classflow_api",
		Subsystem: "http",
		Name:      "schedule_request_total",
		Help:      "Total HTTP requests by path and group",
	},
	[]string{"endpoint", "group_name"},
)

func init() {
	prometheus.MustRegister(httpRequestsCounter)
	prometheus.MustRegister(scheduleRequestCounter)
}

func IncRequestCounter(endpoint string) {
	httpRequestsCounter.WithLabelValues(endpoint).Inc()
}

func IncScheduleRequestCounter(endpoint, groupName string) {
	scheduleRequestCounter.WithLabelValues(endpoint, groupName).Inc()
}
