package messagehandler

import (
	"encoding/json"
	"log/slog"

	"codingminds.com/homemmonitor/metrics"
	"codingminds.com/homemmonitor/model"
)

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
	metrics.ObsCurrentPowerConsumption(msg.Payload.Data.LiveMeasurement.Power)
	metrics.ObsAveragePower(msg.Payload.Data.LiveMeasurement.AveragePower)
	metrics.ObsAccumulatedConsumption(msg.Payload.Data.LiveMeasurement.AccumulatedConsumption)
	metrics.ObsAccumulatedCost(msg.Payload.Data.LiveMeasurement.AccumulatedCost)
	metrics.ObsMinPower(msg.Payload.Data.LiveMeasurement.MinPower)
	metrics.ObsMaxPower(msg.Payload.Data.LiveMeasurement.MaxPower)
	metrics.ObsPhaseCurrent("1", msg.Payload.Data.LiveMeasurement.CurrentL1)
	metrics.ObsPhaseCurrent("2", msg.Payload.Data.LiveMeasurement.CurrentL2)
	metrics.ObsPhaseCurrent("3", msg.Payload.Data.LiveMeasurement.CurrentL3)
	metrics.ObsPhaseVoltage("1", msg.Payload.Data.LiveMeasurement.VoltagePhase1)
	metrics.ObsPhaseVoltage("2", msg.Payload.Data.LiveMeasurement.VoltagePhase2)
	metrics.ObsPhaseVoltage("3", msg.Payload.Data.LiveMeasurement.VoltagePhase3)
}

func logMessage(msg model.LiveMeasurementResponseBody) {
	slog.Info("Measurements received from Tibber",
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
