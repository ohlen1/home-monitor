package metrics

import (
	"log/slog"
	"net/http"

	"codingminds.com/homemmonitor/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var currenPowerGauge *prometheus.GaugeVec = nil
var currentPhaseCurrent *prometheus.GaugeVec = nil
var currentPhaseVoltage *prometheus.GaugeVec = nil
var averagePower *prometheus.GaugeVec = nil
var minPower *prometheus.GaugeVec = nil
var maxPower *prometheus.GaugeVec = nil
var accumulatedConsumption *prometheus.GaugeVec = nil
var accumulatedCost *prometheus.GaugeVec = nil

var cfg config.Config

func Init(config config.Config) {
	cfg = config
	if cfg.MetricsEnabled {
		currenPowerGauge = newGaugeVector("hm_current_power_gauge", "Gauge for current power consumption", []string{"homeId"})
		currentPhaseCurrent = newGaugeVector("hm_current_phase_current", "Gauge for current phase current", []string{"homeId", "phaseNo"})
		currentPhaseVoltage = newGaugeVector("hm_current_phase_voltage", "Gauge for current phase voltage", []string{"homeId", "phaseNo"})
		averagePower = newGaugeVector("hm_average_power_since_midnight", "Gauge for average power consumption, since midnight", []string{"homeId"})
		minPower = newGaugeVector("hm_min_power_since_midnight", "Gauge for minimum power consumption, since midnight", []string{"homeId"})
		maxPower = newGaugeVector("hm_maxpower_since_midnight", "Gauge for maximum power consumption, since midnight", []string{"homeId"})
		accumulatedConsumption = newGaugeVector("hm_accumulated_consumption_since_midnight", "Gauge for accumulated consumption in SEK, since midnight", []string{"homeId"})
		accumulatedCost = newGaugeVector("hm_accumulated_cost_since_midnight", "Gauge for accumulated cost in SEK, since midnight", []string{"homeId"})
		go Listen()
		slog.Info("Metrics initialized")
	}
}

func newGauge(name, help string) prometheus.Gauge {
	o := prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}
	gauge := prometheus.NewGauge(o)
	prometheus.MustRegister(gauge)
	return gauge
}

func newGaugeVector(name, help string, labels []string) *prometheus.GaugeVec {
	gaugeOpts := prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}

	gauge := prometheus.NewGaugeVec(gaugeOpts, labels)
	prometheus.MustRegister(gauge)
	return gauge
}

func newCounterVector(name, help string, labels []string) *prometheus.CounterVec {
	counterOpts := prometheus.CounterOpts{
		Name: name,
		Help: help,
	}
	counter := prometheus.NewCounterVec(counterOpts, labels)
	prometheus.MustRegister(counter)
	return counter
}

func ObsCurrentPowerConsumption(v float64) {
	if currenPowerGauge != nil {
		currenPowerGauge.WithLabelValues(cfg.Tibber.HomeId).Add(v)
	}
}

func ObsPhaseCurrent(phaseNo string, v float64) {
	if currentPhaseCurrent != nil {
		currentPhaseCurrent.WithLabelValues(cfg.Tibber.HomeId, phaseNo).Add(v)
	}
}

func ObsPhaseVoltage(phaseNo string, v float64) {
	if currentPhaseVoltage != nil {
		currentPhaseVoltage.WithLabelValues(cfg.Tibber.HomeId, phaseNo).Add(v)
	}
}

func ObsAveragePower(watts float64) {
	if averagePower != nil {
		averagePower.WithLabelValues(cfg.Tibber.HomeId).Add(watts)
	}
}

func ObsMinPower(watts float64) {
	if minPower != nil {
		minPower.WithLabelValues(cfg.Tibber.HomeId).Add(watts)
	}
}

func ObsMaxPower(watts float64) {
	if maxPower != nil {
		maxPower.WithLabelValues(cfg.Tibber.HomeId).Add(watts)
	}
}

func ObsAccumulatedConsumption(cost float64) {
	if accumulatedCost != nil {
		accumulatedCost.WithLabelValues(cfg.Tibber.HomeId).Add(cost)
	}
}

func ObsAccumulatedCost(cost float64) {
	if accumulatedCost != nil {
		accumulatedCost.WithLabelValues(cfg.Tibber.HomeId).Add(cost)
	}
}

func AddHandler() {
	http.Handle("/metrics", promhttp.Handler())
	slog.Info("Added prometheus handler", "endpoint", "/metrics")
}

func Listen() {
	slog.Info("Metrics server enabled on port 7071")
	AddHandler()
	http.ListenAndServe(":7071", nil)
}
