package measurements

import (
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
