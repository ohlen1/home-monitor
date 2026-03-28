package messagehandler

import (
	"encoding/json"
	"log/slog"

	"codingminds.com/homemmonitor/metrics"
	"codingminds.com/homemmonitor/model"
)

// MetricsObserver abstracts the metrics package so it can be replaced in tests.
type MetricsObserver interface {
	ObsCurrentPowerConsumption(v float64)
	ObsAveragePower(watts float64)
	ObsAccumulatedConsumption(kwh float64)
	ObsAccumulatedCost(cost float64)
	ObsMinPower(watts float64)
	ObsMaxPower(watts float64)
	ObsPhaseCurrent(phaseNo string, v float64)
	ObsPhaseVoltage(phaseNo string, v float64)
}

var obs MetricsObserver = realMetrics{}

type realMetrics struct{}

func (realMetrics) ObsCurrentPowerConsumption(v float64)      { metrics.ObsCurrentPowerConsumption(v) }
func (realMetrics) ObsAveragePower(watts float64)             { metrics.ObsAveragePower(watts) }
func (realMetrics) ObsAccumulatedConsumption(kwh float64)     { metrics.ObsAccumulatedConsumption(kwh) }
func (realMetrics) ObsAccumulatedCost(cost float64)           { metrics.ObsAccumulatedCost(cost) }
func (realMetrics) ObsMinPower(watts float64)                 { metrics.ObsMinPower(watts) }
func (realMetrics) ObsMaxPower(watts float64)                 { metrics.ObsMaxPower(watts) }
func (realMetrics) ObsPhaseCurrent(phaseNo string, v float64) { metrics.ObsPhaseCurrent(phaseNo, v) }
func (realMetrics) ObsPhaseVoltage(phaseNo string, v float64) { metrics.ObsPhaseVoltage(phaseNo, v) }

func Handle(data []byte) {
	msg := new(model.LiveMeasurementResponseBody)
	err := json.Unmarshal(data, msg)
	if err != nil {
		slog.Error("Failed to unmarshal LiveMeasurements", "error", err)
		return
	}
	if msg.Id != "1" || msg.Type != "next" {
		slog.Warn("Message Id and Type has unexpected values", "id", msg.Id, "type", msg.Type)
	}

	logMessage(*msg)
	produceMetrics(*msg)
}

func produceMetrics(msg model.LiveMeasurementResponseBody) {
	obs.ObsCurrentPowerConsumption(msg.Payload.Data.LiveMeasurement.Power)
	obs.ObsAveragePower(msg.Payload.Data.LiveMeasurement.AveragePower)
	obs.ObsAccumulatedConsumption(msg.Payload.Data.LiveMeasurement.AccumulatedConsumption)
	obs.ObsAccumulatedCost(msg.Payload.Data.LiveMeasurement.AccumulatedCost)
	obs.ObsMinPower(msg.Payload.Data.LiveMeasurement.MinPower)
	obs.ObsMaxPower(msg.Payload.Data.LiveMeasurement.MaxPower)
	obs.ObsPhaseCurrent("1", msg.Payload.Data.LiveMeasurement.CurrentL1)
	obs.ObsPhaseCurrent("2", msg.Payload.Data.LiveMeasurement.CurrentL2)
	obs.ObsPhaseCurrent("3", msg.Payload.Data.LiveMeasurement.CurrentL3)
	obs.ObsPhaseVoltage("1", msg.Payload.Data.LiveMeasurement.VoltagePhase1)
	obs.ObsPhaseVoltage("2", msg.Payload.Data.LiveMeasurement.VoltagePhase2)
	obs.ObsPhaseVoltage("3", msg.Payload.Data.LiveMeasurement.VoltagePhase3)
}

func logMessage(msg model.LiveMeasurementResponseBody) {
	slog.Debug("Measurements received from Tibber",
		"timestamp", msg.Payload.Data.LiveMeasurement.Timestamp,
		"power", msg.Payload.Data.LiveMeasurement.Power,
		"averagePower", msg.Payload.Data.LiveMeasurement.AveragePower,
		"accumulatedConsumption", msg.Payload.Data.LiveMeasurement.AccumulatedConsumption,
		"accumulatedCost", msg.Payload.Data.LiveMeasurement.AccumulatedCost,
		"minPower", msg.Payload.Data.LiveMeasurement.MinPower,
		"maxPower", msg.Payload.Data.LiveMeasurement.MaxPower,
		"currentL1", msg.Payload.Data.LiveMeasurement.CurrentL1,
		"currentL2", msg.Payload.Data.LiveMeasurement.CurrentL2,
		"currentL3", msg.Payload.Data.LiveMeasurement.CurrentL3,
		"voltagePhase1", msg.Payload.Data.LiveMeasurement.VoltagePhase1,
		"voltagePhase2", msg.Payload.Data.LiveMeasurement.VoltagePhase2,
		"voltagePhase3", msg.Payload.Data.LiveMeasurement.VoltagePhase3)
}
