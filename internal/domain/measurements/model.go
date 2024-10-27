package measurements

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

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

func (m *Measurement) CalculateStateOfEnergy(max Power) {
	// To prevent division by zero
	if max.Value == 0 {
		m.StateOfEnergy = 0
		return
	}

	m.StateOfEnergy = math.Abs(m.Power.Value/max.Value) * 100
}

func NewRandomMeasurement(min, max Power, powerStep float64) Measurement {
	// Todo validate min and max values
	// Todo validate powerStep value

	// Generate a random power value between min and max, with a step of powerStep

	// If powerStep is 0, there is no limit on power step
	if powerStep <= 0 {

	}

	rangeSize := int((max.Value - min.Value) / powerStep)
	randomStep, _ := rand.Int(rand.Reader, big.NewInt(int64(rangeSize+1)))
	randomValue := min.Value + float64(randomStep.Int64())*powerStep

	measurement := Measurement{
		Power: Power{
			Value: randomValue,
			Unit:  UnitWatt,
		},
		Time: time.Now(),
	}

	measurement.CalculateStateOfEnergy(max)
	return measurement
}
