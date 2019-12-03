package mapper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tendermint/tendermint/config"
)

const (
	MetricsSubsystem = "mint" // module name
)

type Metrics struct {
	TotalAppliedQOS   metrics.Gauge // total applied QOS
	MintPerBlockQOS   metrics.Gauge // mint QOS per block
	GasFeePerBlockQOS metrics.Gauge // gas fee per block
}

func NopMetrics() *Metrics {
	return &Metrics{
		TotalAppliedQOS:   discard.NewGauge(),
		MintPerBlockQOS:   discard.NewGauge(),
		GasFeePerBlockQOS: discard.NewGauge(),
	}
}

// 注册本模块监控项
func PrometheusMetrics(cfg *config.InstrumentationConfig) *Metrics {
	if !cfg.Prometheus {
		return NopMetrics()
	}

	appliedQOSVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "total_applied_qos",
		Help:      "total applied qos",
	}, []string{})
	mintQOSVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "mint_per_block",
		Help:      "mint qos per block",
	}, []string{})
	gasFeeVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "gas_fee_per_block",
		Help:      "gas fee per block",
	}, []string{})

	err := stdprometheus.Register(appliedQOSVec)
	if err != nil {
		panic(err)
	}
	err = stdprometheus.Register(mintQOSVec)
	if err != nil {
		panic(err)
	}
	err = stdprometheus.Register(gasFeeVec)
	if err != nil {
		panic(err)
	}

	return &Metrics{
		TotalAppliedQOS:   prometheus.NewGauge(appliedQOSVec),
		MintPerBlockQOS:   prometheus.NewGauge(mintQOSVec),
		GasFeePerBlockQOS: prometheus.NewGauge(gasFeeVec),
	}
}
