package mapper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tendermint/tendermint/config"
)

const (
	MetricsSubsystem = "gov"
	VoterLabel       = "voter"
	ProposalIdLabel  = "proposal"
	OptionLabel      = "option"
)

type Metrics struct {
	ProposalStatus metrics.Gauge // StatusNil:0 StatusDepositPeriod:1 StatusVotingPeriod:2 StatusPassed:3 StatusRejected:4
	Vote           metrics.Gauge // Total power of the voting delegations
}

func NopMetrics() *Metrics {
	return &Metrics{
		ProposalStatus: discard.NewGauge(),
		Vote:           discard.NewGauge(),
	}
}

func PrometheusMetrics(cfg *config.InstrumentationConfig) *Metrics {
	if !cfg.Prometheus {
		return NopMetrics()
	}

	statusVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "proposal_status",
		Help:      "proposal status",
	}, []string{ProposalIdLabel})
	voteVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "vote",
		Help:      "total power of the voting delegations",
	}, []string{ProposalIdLabel, VoterLabel, OptionLabel})
	err := stdprometheus.Register(statusVec)
	if err != nil {
		panic(err)
	}
	err = stdprometheus.Register(voteVec)
	if err != nil {
		panic(err)
	}

	return &Metrics{
		ProposalStatus: prometheus.NewGauge(statusVec),
		Vote:           prometheus.NewGauge(voteVec),
	}
}
