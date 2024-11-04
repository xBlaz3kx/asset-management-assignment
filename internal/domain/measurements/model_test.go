package measurements

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateStateOfEnergy(t *testing.T) {
	tests := []struct {
		name                string
		measurement         Measurement
		previousMeasurement Measurement
		expected            float32
	}{
		{
			name: "20% energy",
			measurement: Measurement{
				Power: Power{
					Value: 20.0,
					Unit:  UnitWatt,
				},
			},
			previousMeasurement: Measurement{
				Power: Power{
					Value: 100.0,
				},
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
			previousMeasurement: Measurement{
				Power: Power{
					Value: 100.0,
				},
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
			previousMeasurement: Measurement{
				Power: Power{
					Value: 100.0,
				},
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
			previousMeasurement: Measurement{
				Power: Power{
					Value: 100.0,
				},
			},
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.measurement.CalculateStateOfEnergy(tt.previousMeasurement)
			assert.InDelta(t, tt.expected, tt.measurement.StateOfEnergy, 0.001)
		})
	}
}
