package config

import (
	"fmt"
	"slices"
	"testing"
)

var tibberMeasurements []string = []string{"timestamp", "power", "accumulatedConsumption", "accumulatedCost", "currency", "minPower", "averagePower", "maxPower"}

func TestConfig(t *testing.T) {
	cfg := LoadConfig("../test/config.yaml")
	if cfg.ServiceName != "home-monitor" {
		t.Errorf("Failed ServiceName configuration, should be test but was: %s", cfg.ServiceName)
	}
	if cfg.Parallelism != 5 {
		t.Errorf("Failed Parallelism configuration, should be 5 but was: %s", fmt.Sprint(cfg.Parallelism))
	}

	if !slices.Contains(cfg.Tibber.Measurements, "timestamp") {
		t.Errorf("tibber.measurements does not contain timestamp")
	}

	if !validTibberMeasurements(t, cfg.Tibber.Measurements) {
		s := fmt.Sprintf("tibber.measurements does not match. Expected: [%s]. Found: [%s]", sliceToCommaSeparatedString(tibberMeasurements), sliceToCommaSeparatedString(cfg.Tibber.Measurements))
		t.Errorf(s)
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

func validTibberMeasurements(t *testing.T, m []string) bool {
	for _, s := range tibberMeasurements {
		if !slices.Contains(m, s) {
			return false
		}
	}
	return true
}
