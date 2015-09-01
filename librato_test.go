package librato

import (
	"testing"
	"time"

	"github.com/rcrowley/go-metrics"
)

func TestDefaultRateOptions(t *testing.T) {
	r := metrics.DefaultRegistry
	p := NewReporter(
		r,
		time.Second*5, // interval
		"",            // account owner email address
		"",            // Librato API token
		"",            // source
		[]float64{0.99, 0.90, 0.50}, // percentiles to send
		time.Millisecond,            // time unit
	)
	ts := time.Now()
	time.Sleep(5 * time.Millisecond)
	metrics.GetOrRegisterTimer("test", r).UpdateSince(ts)
	now := time.Now()
	b, err := p.BuildRequest(now, r)
	if err != nil {
		t.Error("Librato initialization failed with: %v", err)
	}

	r1, r5, r15 := false, false, false

	for _, g := range b.Gauges {
		for k, v := range g {
			if k == "name" {
				if v == "test.rate.1min" {
					r1 = true
				} else if v == "test.rate.5min" {
					r5 = true
				} else if v == "test.rate.15min" {
					r15 = true
				}
			}
		}
	}

	if !r1 || !r5 || !r15 {
		t.Error("Expected Timer Rate function - but got none")
	}

}

func TestNoRateOptions(t *testing.T) {
	r := metrics.DefaultRegistry
	p := NewReporterWithRateOptions(
		r,
		RateOptions{}, // no rates por favor
		time.Second*5, // interval
		"",            // account owner email address
		"",            // Librato API token
		"",            // source
		[]float64{0.99, 0.90, 0.50}, // percentiles to send
		time.Millisecond,            // time unit
	)
	ts := time.Now()
	time.Sleep(5 * time.Millisecond)
	metrics.GetOrRegisterTimer("test", r).UpdateSince(ts)
	now := time.Now()
	b, err := p.BuildRequest(now, r)
	if err != nil {
		t.Error("Librato initialization failed with: %v", err)
	}

	r1, r5, r15 := false, false, false

	for _, g := range b.Gauges {
		for k, v := range g {
			if k == "name" {
				if v == "test.rate.1min" {
					r1 = true
				} else if v == "test.rate.5min" {
					r5 = true
				} else if v == "test.rate.15min" {
					r15 = true
				}
			}
		}
	}

	if r1 || r5 || r15 {
		t.Error("Expected No Timer Rate function - but got at least one")
	}

}
