package mapper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tendermint/tendermint/config"
)

const (
	MetricsSubsystem = "params"
	ParamLabel       = "key"
)

type Metrics struct {
	Param metrics.Gauge
}

func NopMetrics() *Metrics {
	return &Metrics{
		Param: discard.NewGauge(),
	}
}

func PrometheusMetrics(cfg *config.InstrumentationConfig) *Metrics {
	if !cfg.Prometheus {
		return NopMetrics()
	}

	paramVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "parameter",
		Help:      "parameter changes",
	}, []string{ParamLabel})
	err := stdprometheus.Register(paramVec)
	if err != nil {
		panic(err)
	}

	return &Metrics{
		Param: prometheus.NewGauge(paramVec),
	}
}
