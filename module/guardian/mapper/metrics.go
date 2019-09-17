package mapper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tendermint/tendermint/config"
)

const (
	MetricsSubsystem = "guardian"
	OperationLabel   = "operation"
	AddressLabel     = "address"
)

type Metrics struct {
	Guardian metrics.Gauge // guardian operations
}

func NopMetrics() *Metrics {
	return &Metrics{
		Guardian: discard.NewGauge(),
	}
}

func PrometheusMetrics(cfg *config.InstrumentationConfig) *Metrics {
	if !cfg.Prometheus {
		return NopMetrics()
	}

	guardianVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "guardian_operation",
		Help:      "operations by guardian",
	}, []string{AddressLabel, OperationLabel})

	err := stdprometheus.Register(guardianVec)
	if err != nil {
		panic(err)
	}

	return &Metrics{
		Guardian: prometheus.NewGauge(guardianVec),
	}
}
