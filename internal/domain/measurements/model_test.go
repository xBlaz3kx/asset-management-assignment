package measurements

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateStateOfEnergy(t *testing.T) {
	tests := []struct {
		name        string
		measurement Measurement
		maxPower    Power
		expected    float32
	}{
		{
			name: "20% energy",
			measurement: Measurement{
				Power: Power{
					Value: 20.0,
					Unit:  UnitWatt,
				},
			},
			maxPower: Power{
				Value: 100.0,
			},
			expected: 20.0,
		},
		{
			name: "50% energy",
			measurement: Measurement{
				Power: Power{
					Value: 50.0,
					Unit:  UnitWatt,
				},
			},
			maxPower: Power{
				Value: 100.0,
			},
			expected: 50.0,
		},
		{
			name: "0% energy",
			measurement: Measurement{
				Power: Power{
					Value: 0.0,
					Unit:  UnitWatt,
				},
			},
			maxPower: Power{
				Value: 100.0,
			},
			expected: 0.0,
		},
		{
			name: "100% energy",
			measurement: Measurement{
				Power: Power{
					Value: 100.0,
					Unit:  UnitWatt,
				},
			},
			maxPower: Power{
				Value: 100.0,
			},
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.measurement.CalculateStateOfEnergy(tt.maxPower)
			assert.InDelta(t, tt.expected, tt.measurement.StateOfEnergy, 0.001)
		})
	}
}

func TestNewRandomMeasurement(t *testing.T) {
	tests := []struct {
		name         string
		minPower     Power
		maxPower     Power
		maxPowerStep float64
		expected     Measurement
	}{
		{
			name: "Random measurement",
			minPower: Power{
				Value: 0,
				Unit:  UnitWatt,
			},
			maxPower: Power{
				Value: 100,
				Unit:  UnitWatt,
			},
			maxPowerStep: 10,
		},
		{
			name: "Unlimited power step",
			minPower: Power{
				Value: 0,
				Unit:  UnitWatt,
			},
			maxPower: Power{
				Value: 100,
				Unit:  UnitWatt,
			},
			maxPowerStep: -1,
		},
		{
			name: "Invalid max power",
			minPower: Power{
				Value: 100,
				Unit:  UnitWatt,
			},
			maxPower: Power{
				Value: 0,
				Unit:  UnitWatt,
			},
			maxPowerStep: 10,
		},
		{
			name: "Invalid min power",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewRandomMeasurement(tt.minPower, tt.maxPower, tt.maxPowerStep)
			assert.InDelta(t, tt.minPower.Value, actual.Power.Value, tt.maxPowerStep)
		})
	}
}
