package generator

import (
	"math"
	"math/rand/v2"
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
)

func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func zeroPowerMeasurement(stateOfEnergy float64) *measurements.Measurement {
	return &measurements.Measurement{
		Power: measurements.Power{
			Value: 0,
			Unit:  "W",
		},
		StateOfEnergy: stateOfEnergy,
		Time:          time.Now(),
	}
}

func getPowerStep(maxPowerStep float64) float64 {
	stepValue := maxPowerStep

	// Generate a random step value
	if maxPowerStep <= 0 {
		stepValue = rand.Float64() * maxPowerStep
	}

	sign := rand.IntN(2)
	if sign == 0 {
		// Negative step
		sign = -1
	} else {
		// Positive step
		sign = 1
	}

	return math.Copysign(stepValue, float64(sign))
}
