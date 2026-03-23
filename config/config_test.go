package config

import (
	"slices"
	"testing"
)

var tibberMeasurements = []string{"timestamp", "power", "accumulatedConsumption", "accumulatedCost", "currency", "minPower", "averagePower", "maxPower"}

func TestConfig(t *testing.T) {
	cfg := LoadConfig("../testdata/config-test.yaml")
	if cfg.ServiceName != "home-monitor" {
		t.Fatalf("Failed ServiceName configuration, should be home-monitor but was: %s", cfg.ServiceName)
	}
	if cfg.Parallelism != 5 {
		t.Fatalf("Failed Parallelism configuration, should be 5 but was: %d", cfg.Parallelism)
	}

	if !slices.Contains(cfg.Tibber.Measurements, "timestamp") {
		t.Fatalf("tibber.measurements does not contain timestamp")
	}

	if !validTibberMeasurements(cfg.Tibber.Measurements) {
		t.Fatalf("tibber.measurements does not match. Expected: [%s]. Found: [%s]", sliceToCommaSeparatedString(tibberMeasurements), sliceToCommaSeparatedString(cfg.Tibber.Measurements))
	}
}

func sliceToCommaSeparatedString(s []string) string {
	msg := ""
	for i, a := range s {
		msg += a
		if i < len(s)-1 {
			msg += ", "
		}
	}
	return msg
}

func validTibberMeasurements(m []string) bool {
	for _, s := range tibberMeasurements {
		if !slices.Contains(m, s) {
			return false
		}
	}
	return true
}
