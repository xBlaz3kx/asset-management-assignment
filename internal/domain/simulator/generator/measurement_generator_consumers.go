package generator

import (
	"time"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/measurements"
	"asset-measurements-assignment/internal/domain/simulator"
)

type ConsumerMeasurementGenerator struct {
	cfg simulator.Configuration
	// Last generated measurement
	previousMeasurement *measurements.Measurement
}

func NewConsumer(cfg simulator.Configuration) *ConsumerMeasurementGenerator {
	return &ConsumerMeasurementGenerator{
		cfg: cfg,
	}
}

// GenerateMeasurement generates a random measurement for the asset based on the provided configuration
// and the previous measurement.
func (c *ConsumerMeasurementGenerator) GenerateMeasurement() (*measurements.Measurement, error) {
	if c.previousMeasurement == nil {
		// This is the first measurement
		measurement := zeroPowerMeasurement(0)
		c.previousMeasurement = measurement
		return measurement, nil
	}

	powerStep := getPowerStep(c.cfg.MaxPowerStep)
	power := c.previousMeasurement.Power.Value + powerStep
	power = clamp(power, c.cfg.MinPower, c.cfg.MaxPower)

	measurement := &measurements.Measurement{
		Power: measurements.Power{
			Value: power,
			Unit:  measurements.UnitWatt,
		},
		Time: time.Now(),
	}

	c.calculateSoE(measurement)

	c.previousMeasurement = measurement
	return measurement, nil
}

func (c *ConsumerMeasurementGenerator) GetEnergyType() domain.EnergyType {
	return domain.EnergyTypeConsumer
}

func (c *ConsumerMeasurementGenerator) calculateSoE(measurement *measurements.Measurement) {
	// Calculate energy between two measurements
	timeElapsed := measurement.Time.Sub(c.previousMeasurement.Time).Seconds()
	energyChange := ((measurement.Power.Value + c.previousMeasurement.Power.Value) / 2) * timeElapsed
	energyPercentageChange := (energyChange / c.cfg.MaxPower) * 100
	measurement.StateOfEnergy = c.previousMeasurement.StateOfEnergy + energyPercentageChange

	measurement.StateOfEnergy = clamp(measurement.StateOfEnergy, 0, 100)
}