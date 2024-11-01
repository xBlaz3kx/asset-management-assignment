package measurements

import (
	"math"
	"math/rand/v2"
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
func (m *Measurement) CalculateStateOfEnergy(max Power) {
	// To prevent division by zero
	if max.Value == 0 {
		m.StateOfEnergy = 0
		return
	}

	m.StateOfEnergy = math.Abs(m.Power.Value/max.Value) * 100
}

// NewRandomMeasurement generates a random measurement based on the given min and max power values and max power step.
// Min power value must be less than or equal to max power value in order to generate a valid measurement.
// If max power step is less than or equal to 0, the power value will be a random value between min and max power values.
func NewRandomMeasurement(min, max Power, maxPowerStep float64) (*Measurement, error) {
	// Assuming the units are the same
	if min.Value > max.Value {
		return nil, ErrMinPowerGreaterThanMaxPower
	}

	power := 0.0
	if maxPowerStep <= 0.0 {
		power = rand.Float64()*(max.Value-min.Value) + min.Value
	} else {
		rangeSize := (max.Value - min.Value) / maxPowerStep
		randomStep := rand.IntN(int(rangeSize + 1))
		power = min.Value + float64(randomStep)*maxPowerStep
	}

	measurement := Measurement{
		Power: Power{
			Value: power,
			Unit:  UnitWatt,
		},
		Time: time.Now(),
	}

	measurement.CalculateStateOfEnergy(max)
	return &measurement, nil
}
