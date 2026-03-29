package messagehandler

import (
	"testing"
)

type mockMetrics struct {
	currentPower           float64
	currentPowerProduction float64
	averagePower           float64
	accumulatedConsumption float64
	accumulatedCost        float64
	minPower               float64
	maxPower               float64
	phaseCurrent           map[string]float64
	phaseVoltage           map[string]float64
}

func (m *mockMetrics) ObsCurrentPowerConsumption(v float64)      { m.currentPower = v }
func (m *mockMetrics) ObsCurrentPowerProduction(watts float64)   { m.currentPowerProduction = watts }
func (m *mockMetrics) ObsAveragePower(watts float64)             { m.averagePower = watts }
func (m *mockMetrics) ObsAccumulatedConsumption(kwh float64)     { m.accumulatedConsumption = kwh }
func (m *mockMetrics) ObsAccumulatedCost(cost float64)           { m.accumulatedCost = cost }
func (m *mockMetrics) ObsMinPower(watts float64)                 { m.minPower = watts }
func (m *mockMetrics) ObsMaxPower(watts float64)                 { m.maxPower = watts }
func (m *mockMetrics) ObsPhaseCurrent(phaseNo string, v float64) { m.phaseCurrent[phaseNo] = v }
func (m *mockMetrics) ObsPhaseVoltage(phaseNo string, v float64) { m.phaseVoltage[phaseNo] = v }

func TestHandle(t *testing.T) {
	mock := &mockMetrics{
		phaseCurrent: make(map[string]float64),
		phaseVoltage: make(map[string]float64),
	}
	obs = mock
	t.Cleanup(func() { obs = realMetrics{} })

	Handle([]byte(`{"id":"1","type":"next","payload":{"data":{"liveMeasurement":{"timestamp":"2024-06-01T12:00:00Z","power":100.0,"averagePower":50.0,"accumulatedConsumption":10.0,"accumulatedCost":2.5,"minPower":20.0,"maxPower":150.0,"currentL1":1.0,"currentL2":2.0,"currentL3":3.0,"voltagePhase1":230.0,"voltagePhase2":231.0,"voltagePhase3":229.0}}}}`))

	assertFloat(t, "currentPower", 100.0, mock.currentPower)
	assertFloat(t, "currentPowerProduction", 50.0, mock.currentPowerProduction)
	assertFloat(t, "averagePower", 50.0, mock.averagePower)
	assertFloat(t, "accumulatedConsumption", 10.0, mock.accumulatedConsumption)
	assertFloat(t, "accumulatedCost", 2.5, mock.accumulatedCost)
	assertFloat(t, "minPower", 20.0, mock.minPower)
	assertFloat(t, "maxPower", 150.0, mock.maxPower)
	assertFloat(t, "currentL1", 1.0, mock.phaseCurrent["1"])
	assertFloat(t, "currentL2", 2.0, mock.phaseCurrent["2"])
	assertFloat(t, "currentL3", 3.0, mock.phaseCurrent["3"])
	assertFloat(t, "voltagePhase1", 230.0, mock.phaseVoltage["1"])
	assertFloat(t, "voltagePhase2", 231.0, mock.phaseVoltage["2"])
	assertFloat(t, "voltagePhase3", 229.0, mock.phaseVoltage["3"])
}

func assertFloat(t *testing.T, name string, expected, got float64) {
	t.Helper()
	if got != expected {
		t.Errorf("%s: expected %f, got %f", name, expected, got)
	}
}
