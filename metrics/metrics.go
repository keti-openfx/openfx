package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricOptions struct {
	GatewayFunctionInvocation *prometheus.CounterVec
	GatewayFunctionsHistogram *prometheus.HistogramVec
	ServiceReplicasGauge      *prometheus.GaugeVec
}

// BuildMetricsOptions builds metrics for tracking functions in the API gateway
func BuildMetricsOptions() MetricOptions {
	gatewayFunctionsHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "gateway_functions_seconds",
		Help: "Function time taken",
	}, []string{"function_name"})

	gatewayFunctionInvocation := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_function_invocation_total",
			Help: "Individual function metrics",
		},
		[]string{"function_name", "code"})

	serviceReplicas := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_service_count",
			Help: "kubernetes service replicas",
		},
		[]string{"function_name"},
	)

	metricsOptions := MetricOptions{
		GatewayFunctionsHistogram: gatewayFunctionsHistogram,
		GatewayFunctionInvocation: gatewayFunctionInvocation,
		ServiceReplicasGauge:      serviceReplicas,
	}

	return metricsOptions
}

func (m *MetricOptions) Notify(serviceName string, duration time.Duration, code string) {
	seconds := duration.Seconds()
	m.GatewayFunctionsHistogram.WithLabelValues(serviceName).Observe(seconds)
	m.GatewayFunctionInvocation.With(prometheus.Labels{"function_name": serviceName, "code": "200"}).Inc()
}
