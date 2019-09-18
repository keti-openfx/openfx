package metrics

import (
	"log"
	"net/http"
	"time"

	"github.com/keti-openfx/openfx/pb"
	"github.com/keti-openfx/openfx/cmd"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
)

// Exporter is a prometheus exporter
type Exporter struct {
	metricOptions MetricOptions
	services      []*pb.Function
}

// NewExporter creates a new exporter for the gateway metrics
func NewExporter(options MetricOptions) *Exporter {
	return &Exporter{
		metricOptions: options,
		services:      []*pb.Function{},
	}
}

//RegisterMetrics registers with Prometheus for tracking
func RegisterExporter(exporter *Exporter) {
	prometheus.MustRegister(exporter)
}

// PrometheusHandler Bootstraps prometheus for metrics collection
func PrometheusHandler() http.Handler {
	return prometheus.Handler()
}

// Describe is to describe the metrics for Prometheus
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.metricOptions.GatewayFunctionInvocation.Describe(ch)
	e.metricOptions.GatewayFunctionsHistogram.Describe(ch)
	e.metricOptions.ServiceReplicasGauge.Describe(ch)
}

// Collect collects data to be consumed by prometheus
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.metricOptions.GatewayFunctionInvocation.Collect(ch)
	e.metricOptions.GatewayFunctionsHistogram.Collect(ch)

	e.metricOptions.ServiceReplicasGauge.Reset()
	for _, service := range e.services {
		e.metricOptions.ServiceReplicasGauge.
			WithLabelValues(service.Name).
			Set(float64(service.Replicas))
	}
	e.metricOptions.ServiceReplicasGauge.Collect(ch)
}

// StartServiceWatcher starts a ticker and collects service replica counts to expose to prometheus
func (e *Exporter) StartServiceWatcher( functionNamespace string, 
					clientset *kubernetes.Clientset, 
					metricsOptions MetricOptions, 
					interval time.Duration) {
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:

				services, err := cmd.List(functionNamespace, clientset)
				if err != nil {
					log.Println(err)
					continue
				}

				e.services = services
				break

			case <-quit:
				return
			}
		}
	}()
}
