package generator

import (
	"math"
	"testing"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/stretchr/testify/assert"
)

func TestClamp(t *testing.T) {
	tests := []struct {
		name  string
		min   float64
		max   float64
		value float64
		want  float64
	}{
		{
			name:  "Test clamp value is less than min",
			min:   0,
			max:   100,
			value: -1,
			want:  0,
		},
		{
			name:  "Test clamp value is greater than max",
			min:   0,
			max:   100,
			value: 101,
			want:  100,
		},
		{
			name:  "Test clamp value is between min and max",
			min:   0,
			max:   100,
			value: 50,
			want:  50,
		},
		{
			name:  "Test clamp value is equal to min",
			min:   0,
			max:   100,
			value: 0,
			want:  0,
		},
		{
			name:  "Test clamp value is equal to max",
			min:   0,
			max:   100,
			value: 100,
			want:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clamp(tt.value, tt.min, tt.max); got != tt.want {
				t.Errorf("clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroPowerMeasurement(t *testing.T) {
	tests := []struct {
		name          string
		stateOfCharge float64
	}{
		{
			name:          "State of charge is zero",
			stateOfCharge: 0,
		},
		{
			name:          "State of charge is not zero",
			stateOfCharge: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			measurement := zeroPowerMeasurement(tt.stateOfCharge)
			assert.NotNil(t, measurement)
			assert.Equal(t, 0.0, measurement.Power.Value)
			assert.Equal(t, measurements.UnitWatt, measurement.Power.Unit)
			assert.Equal(t, tt.stateOfCharge, measurement.StateOfEnergy)
		})
	}
}

func TestGetPowerStep(t *testing.T) {
	tests := []struct {
		name         string
		maxPowerStep float64
	}{
		{
			name:         "Max power step is zero",
			maxPowerStep: 0,
		},
		{
			name:         "Max power step is positive",
			maxPowerStep: 10,
		},
		{
			name:         "Max power step is negative",
			maxPowerStep: -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := getPowerStep(tt.maxPowerStep)

			if tt.maxPowerStep <= 0 {
				assert.LessOrEqual(t, tt.maxPowerStep, math.Abs(step))
			} else {
				assert.Equal(t, tt.maxPowerStep, math.Abs(step))
			}
		})
	}
}
