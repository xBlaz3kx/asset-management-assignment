package generator

import (
	"math"
	"time"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/measurements"
	"asset-measurements-assignment/internal/domain/simulator"
)

type ProducerMeasurementGenerator struct {
	cfg simulator.Configuration
	// Last generated measurement
	previousMeasurement *measurements.Measurement
}

func NewProducer(cfg simulator.Configuration) *ProducerMeasurementGenerator {
	return &ProducerMeasurementGenerator{
		cfg: cfg,
	}
}

// GenerateMeasurement generates a random measurement for the asset based on the provided configuration
// and the previous measurement.
func (p *ProducerMeasurementGenerator) GenerateMeasurement() (*measurements.Measurement, error) {
	if p.previousMeasurement == nil {
		// This is the first measurement
		measurement := zeroPowerMeasurement(0)
		p.previousMeasurement = measurement
		return measurement, nil
	}

	powerStep := getPowerStep(p.cfg.MaxPowerStep)
	power := p.previousMeasurement.Power.Value + powerStep

	// Assuming that min and max have negative sign, so we need to swap them
	if power < p.cfg.MaxPower {
		power = p.cfg.MaxPower
	} else if power > p.cfg.MinPower {
		power = p.cfg.MinPower
	}

	measurement := &measurements.Measurement{
		Power: measurements.Power{
			Value: power,
			Unit:  measurements.UnitWatt,
		},
		Time: time.Now(),
	}

	p.calculateSoE(measurement)

	p.previousMeasurement = measurement
	return measurement, nil
}

func (p *ProducerMeasurementGenerator) GetEnergyType() domain.EnergyType {
	return domain.EnergyTypeProducer
}

func (p *ProducerMeasurementGenerator) calculateSoE(measurement *measurements.Measurement) {
	// Calculate energy between two measurements
	timeElapsed := measurement.Time.Sub(p.previousMeasurement.Time).Seconds()
	energyChange := ((measurement.Power.Value + p.previousMeasurement.Power.Value) / 2) * timeElapsed
	energyPercentageChange := (energyChange / math.Abs(p.cfg.MaxPower)) * 100
	measurement.StateOfEnergy = p.previousMeasurement.StateOfEnergy + energyPercentageChange

	measurement.StateOfEnergy = clamp(measurement.StateOfEnergy, 0, 100)
}
