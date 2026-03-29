package metrics

import (
	"testing"

	"codingminds.com/homemmonitor/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestInit(t *testing.T) {
	config := config.LoadConfig("../testdata/config-test.yaml")
	Init(*config)

	validateNotNilGauge(t, currentPowerGauge, "currentPowerGauge")
	validateNotNilGauge(t, currentPowerProductionGauge, "currentPowerProductionGauge")
	validateNotNilGauge(t, currentPhaseCurrent, "currentPhaseCurrent")
	validateNotNilGauge(t, currentPhaseVoltage, "currentPhaseVoltage")
	validateNotNilGauge(t, averagePower, "averagePower")
	validateNotNilGauge(t, minPower, "minPower")
	validateNotNilGauge(t, maxPower, "maxPower")
	validateNotNilGauge(t, minPowerProduction, "minPowerProduction")
	validateNotNilGauge(t, maxPowerProduction, "maxPowerProduction")
	validateNotNilGauge(t, lastMeterProduction, "lastMeterProduction")
	validateNotNilGauge(t, accumulatedConsumption, "accumulatedConsumption")
	validateNotNilGauge(t, accumulatedCost, "accumulatedCost")
	validateNotNilGauge(t, accumulatedProduction, "accumulatedProduction")
}

func TestObsCurrentPowerConsumption(t *testing.T) {
	const expected = 100.0
	ObsCurrentPowerConsumption(expected)

	testGaugeValue(t, currentPowerGauge, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsCurrentPowerProduction(t *testing.T) {
	const expected = 50.0
	ObsCurrentPowerProduction(expected)

	testGaugeValue(t, currentPowerProductionGauge, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsPhaseCurrent(t *testing.T) {
	const expected = 10.0
	ObsPhaseCurrent("1", expected)

	testGaugeValue(t, currentPhaseCurrent, []string{cfg.Tibber.HomeId, "1"}, expected)
}

func TestObsPhaseVoltage(t *testing.T) {
	const expected = 230.0
	ObsPhaseVoltage("1", expected)

	testGaugeValue(t, currentPhaseVoltage, []string{cfg.Tibber.HomeId, "1"}, expected)
}

func TestObsAveragePower(t *testing.T) {
	const expected = 50.0
	ObsAveragePower(expected)

	testGaugeValue(t, averagePower, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsMinPower(t *testing.T) {
	const expected = 20.0
	ObsMinPower(expected)

	testGaugeValue(t, minPower, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsMaxPower(t *testing.T) {
	const expected = 150.0
	ObsMaxPower(expected)

	testGaugeValue(t, maxPower, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsAccumulatedConsumption(t *testing.T) {
	const expected = 5.0
	ObsAccumulatedConsumption(expected)

	testGaugeValue(t, accumulatedConsumption, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsAccumulatedCost(t *testing.T) {
	const expected = 25.0
	ObsAccumulatedCost(expected)

	testGaugeValue(t, accumulatedCost, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsAccumulatedProduction(t *testing.T) {
	const expected = 10.0
	ObsAccumulatedProduction(expected)

	testGaugeValue(t, accumulatedProduction, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsMinPowerProduction(t *testing.T) {
	const expected = 15.0
	ObsMinPowerProduction(expected)

	testGaugeValue(t, minPowerProduction, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsMaxPowerProduction(t *testing.T) {
	const expected = 30.0
	ObsMaxPowerProduction(expected)

	testGaugeValue(t, maxPowerProduction, []string{cfg.Tibber.HomeId}, expected)
}

func TestObsLastMeterProduction(t *testing.T) {
	const expected = 40.0
	ObsLastMeterProduction(expected)

	testGaugeValue(t, lastMeterProduction, []string{cfg.Tibber.HomeId}, expected)
}

func testGaugeValue(t *testing.T, gauge *prometheus.GaugeVec, labels []string, expected float64) {
	metric, err := gauge.GetMetricWithLabelValues(labels...)
	if err != nil {
		t.Fatalf("Failed to get metric with labels %v: %v", labels, err)
	}

	if got := testutil.ToFloat64(metric); got != expected {
		t.Fatalf("Expected gauge value to be %f but got %f for labels %v", expected, got, labels)
	}
}

func validateNotNilGauge(t *testing.T, gauge *prometheus.GaugeVec, name string) {
	if gauge == nil {
		t.Fatalf("%s should not be nil after Init", name)
	}
}
