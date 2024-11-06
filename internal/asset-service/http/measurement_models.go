package http

import (
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
)

// swagger:model
type Measurement struct {
	// swagger:type string
	Timestamp time.Time `json:"timestamp"`

	// Power represents the power of the asset.
	Power Power `json:"power"`

	// StateOfEnergy represents the current state of energy of the asset.
	StateOfEnergy float64 `json:"stateOfEnergy"`
}

// swagger:model
type Power struct {
	// Value represents the value of the power.
	Value float64 `json:"value"`

	// Unit represents the unit of the power.
	Unit string `json:"unit"`
}

// swagger:parameters getMeasurementsAvgWithinTimeInterval
type AssetMeasurementAveragedParams struct {
	TimeRange
	GroupBy string `form:"groupBy" binding:"required,oneof=minute hour 15min"`
	Sort    string `form:"sort" binding:"required,oneof=asc desc"`
}

func (a *AssetMeasurementAveragedParams) toDomainModel() measurements.AssetMeasurementAveragedParams {
	return measurements.AssetMeasurementAveragedParams{
		TimeRange: a.TimeRange.toDomainModel(),
		GroupBy:   a.GroupBy,
		Sort:      a.Sort,
	}
}

// swagger:parameters getMeasurementsWithinTimeInterval
type TimeRange struct {
	From *time.Time `form:"from" binding:"required"`
	To   *time.Time `form:"to" binding:"required"`
}

func (t *TimeRange) toDomainModel() measurements.TimeRange {
	return measurements.TimeRange{
		From: t.From,
		To:   t.To,
	}
}
