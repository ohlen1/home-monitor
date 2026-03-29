package metrics

import (
	"log/slog"
	"net/http"

	"codingminds.com/homemmonitor/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var currentPowerGauge *prometheus.GaugeVec
var currentPowerProductionGauge *prometheus.GaugeVec
var currentPhaseCurrent *prometheus.GaugeVec
var currentPhaseVoltage *prometheus.GaugeVec
var averagePower *prometheus.GaugeVec
var minPower *prometheus.GaugeVec
var maxPower *prometheus.GaugeVec
var minPowerProduction *prometheus.GaugeVec
var maxPowerProduction *prometheus.GaugeVec
var lastMeterProduction *prometheus.GaugeVec
var accumulatedConsumption *prometheus.GaugeVec
var accumulatedCost *prometheus.GaugeVec
var accumulatedProduction *prometheus.GaugeVec

var cfg config.Config

func Init(config config.Config) {
	cfg = config
	if cfg.MetricsEnabled {
		currentPowerGauge = newGaugeVector("hm_current_power_gauge", "Gauge for current power consumption", []string{"homeId"})
		currentPowerProductionGauge = newGaugeVector("hm_current_power_production_gauge", "Gauge for current power production", []string{"homeId"})
		currentPhaseCurrent = newGaugeVector("hm_current_phase_current", "Gauge for current phase current", []string{"homeId", "phaseNo"})
		currentPhaseVoltage = newGaugeVector("hm_current_phase_voltage", "Gauge for current phase voltage", []string{"homeId", "phaseNo"})
		averagePower = newGaugeVector("hm_average_power_since_midnight", "Gauge for average power consumption, since midnight", []string{"homeId"})
		minPower = newGaugeVector("hm_min_power_since_midnight", "Gauge for minimum power consumption, since midnight", []string{"homeId"})
		maxPower = newGaugeVector("hm_maxpower_since_midnight", "Gauge for maximum power consumption, since midnight", []string{"homeId"})
		minPowerProduction = newGaugeVector("hm_min_power_production_since_midnight", "Gauge for minimum power production, since midnight", []string{"homeId"})
		maxPowerProduction = newGaugeVector("hm_max_power_production_since_midnight", "Gauge for maximum power production, since midnight", []string{"homeId"})
		lastMeterProduction = newGaugeVector("hm_last_meter_production_since_midnight", "Gauge for last meter production, since midnight", []string{"homeId"})
		accumulatedConsumption = newGaugeVector("hm_accumulated_consumption_since_midnight", "Gauge for accumulated consumption in kWh, since midnight", []string{"homeId"})
		accumulatedCost = newGaugeVector("hm_accumulated_cost_since_midnight", "Gauge for accumulated cost in SEK, since midnight", []string{"homeId"})
		accumulatedProduction = newGaugeVector("hm_accumulated_production_since_midnight", "Gauge for accumulated production in kWh, since midnight", []string{"homeId"})
		go Listen()
		slog.Info("Metrics initialized")
	}
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

func ObsCurrentPowerConsumption(v float64) {
	setGauge(currentPowerGauge, []string{cfg.Tibber.HomeId}, v)
}

func ObsCurrentPowerProduction(watts float64) {
	setGauge(currentPowerProductionGauge, []string{cfg.Tibber.HomeId}, watts)
}

func ObsPhaseCurrent(phaseNo string, v float64) {
	setGauge(currentPhaseCurrent, []string{cfg.Tibber.HomeId, phaseNo}, v)
}

func ObsPhaseVoltage(phaseNo string, v float64) {
	setGauge(currentPhaseVoltage, []string{cfg.Tibber.HomeId, phaseNo}, v)
}

func ObsAveragePower(watts float64) {
	setGauge(averagePower, []string{cfg.Tibber.HomeId}, watts)
}

func ObsMinPower(watts float64) {
	setGauge(minPower, []string{cfg.Tibber.HomeId}, watts)
}

func ObsMaxPower(watts float64) {
	setGauge(maxPower, []string{cfg.Tibber.HomeId}, watts)
}

func ObsAccumulatedConsumption(kwh float64) {
	setGauge(accumulatedConsumption, []string{cfg.Tibber.HomeId}, kwh)
}

func ObsAccumulatedCost(cost float64) {
	setGauge(accumulatedCost, []string{cfg.Tibber.HomeId}, cost)
}

func ObsMinPowerProduction(watts float64) {
	setGauge(minPowerProduction, []string{cfg.Tibber.HomeId}, watts)
}

func ObsMaxPowerProduction(watts float64) {
	setGauge(maxPowerProduction, []string{cfg.Tibber.HomeId}, watts)
}

func ObsLastMeterProduction(watts float64) {
	setGauge(lastMeterProduction, []string{cfg.Tibber.HomeId}, watts)
}

func ObsAccumulatedProduction(kwh float64) {
	setGauge(accumulatedProduction, []string{cfg.Tibber.HomeId}, kwh)
}

func AddHandler() {
	http.Handle("/metrics", promhttp.Handler())
	slog.Info("Added prometheus handler", "endpoint", "/metrics")
}

func Listen() {
	slog.Info("Metrics server enabled on port 7071")
	AddHandler()
	if err := http.ListenAndServe(":7071", nil); err != nil {
		slog.Error("Metrics server crashed", "error", err)
	}
}

func setGauge(g *prometheus.GaugeVec, labels []string, value float64) {
	if g == nil {
		return
	}
	g.WithLabelValues(labels...).Set(value)
}
