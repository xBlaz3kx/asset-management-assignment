package generator

import (
	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/measurements"
	"asset-measurements-assignment/internal/domain/simulator"
	"github.com/pkg/errors"
)

type MeasurementGenerator interface {
	GenerateMeasurement() (*measurements.Measurement, error)
	GetEnergyType() domain.EnergyType
}

func GetGeneratorFromConfiguration(cfg simulator.Configuration) (MeasurementGenerator, error) {
	switch cfg.Type.GetEnergyType() {
	case domain.EnergyTypeCombined:
		return NewCombined(cfg), nil
	case domain.EnergyTypeConsumer:
		return NewConsumer(cfg), nil
	case domain.EnergyTypeProducer:
		return NewProducer(cfg), nil
	default:
		return nil, errors.New("invalid asset type")
	}
}
