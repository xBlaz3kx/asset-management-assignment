package measurements

import (
	"math"
	"time"

	"github.com/pkg/errors"
)

var ErrMinPowerGreaterThanMaxPower = errors.New("minPower is greater than maxPower")

type Unit string

const UnitWatt Unit = "W"

type Power struct {
	Value float64 `json:"value"`
	Unit  Unit    `json:"unit"`
}

type Measurement struct {
	// Power
	Power Power `json:"power"`

	// In percent
	StateOfEnergy float64 `json:"stateOfEnergy"`

	// Event timestamp
	Time time.Time `json:"time"`
}

// CalculateStateOfEnergy calculates the state of energy based on the given max power value.
func (m *Measurement) CalculateStateOfEnergy(previousMeasurement Measurement) {
	// Calculate energy between two measurements
	energy := (m.Power.Value + previousMeasurement.Power.Value) / 2 * m.Time.Sub(previousMeasurement.Time).Seconds()
	m.StateOfEnergy = math.Abs(previousMeasurement.StateOfEnergy+energy) * 100
}
