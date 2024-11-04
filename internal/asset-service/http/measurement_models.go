package http

import (
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
)

type Measurement struct {
	Timestamp     time.Time `json:"timestamp"`
	Power         Power     `json:"power"`
	StateOfEnergy float64   `json:"stateOfEnergy"`
}

type Power struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

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
