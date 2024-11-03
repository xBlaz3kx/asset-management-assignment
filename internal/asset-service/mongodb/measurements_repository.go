package mongodb

import (
	"context"
	"errors"
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
	mongo2 "asset-measurements-assignment/internal/pkg/infrastructure/mongo"
	"github.com/xBlaz3kx/DevX/observability"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

const collectionName = "asset_measurements"

// Measurement mongo entity
type Measurement struct {
	AssetID       string             `bson:"assetId"`
	Timestamp     time.Time          `bson:"timestamp"`
	Power         measurements.Power `bson:"power"`
	StateOfEnergy float64            `bson:"stateOfEnergy"`
}

type MeasurementsRepository struct {
	obs        observability.Observability
	collection *mongo.Collection
}

func NewMeasurementsRepository(obs observability.Observability, client *mongo.Database) (*MeasurementsRepository, error) {
	tso := options.TimeSeries().
		SetTimeField("timestamp").
		SetMetaField("assetId").
		SetGranularity("seconds")
	opts := options.CreateCollection().SetTimeSeriesOptions(tso)

	// We will just ignore the error for now :)
	_ = client.CreateCollection(context.Background(), collectionName, opts)

	return &MeasurementsRepository{
		obs:        obs,
		collection: client.Collection(collectionName),
	}, nil
}

func (m *MeasurementsRepository) AddMeasurement(ctx context.Context, assetId string, measurement measurements.Measurement) error {
	ctx, cancel := m.obs.Span(ctx, "measurements.repository.AddMeasurement", zap.Any("measurement", measurement))
	defer cancel()

	dbMeasurement := fromMeasurement(assetId, &measurement)
	res, err := m.collection.InsertOne(ctx, dbMeasurement)
	if err != nil {
		return err
	}

	if !res.Acknowledged {
		return errors.New("insert not acknowledged")
	}

	return nil
}

func (m *MeasurementsRepository) GetLatestAssetMeasurement(ctx context.Context, assetID string) (*measurements.Measurement, error) {
	ctx, cancel := m.obs.Span(ctx, "measurements.repository.GetLatestAssetMeasurement", zap.String("assetID", assetID))
	defer cancel()

	// Find the latest measurement for the asset
	opts := options.FindOne()
	opts.SetSort(bson.M{"timestamp": -1})
	res := m.collection.FindOne(ctx, bson.M{"assetId": assetID}, opts)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var dbMeasurement Measurement
	err := res.Decode(&dbMeasurement)
	if err != nil {
		return nil, err
	}

	return toMeasurement(&dbMeasurement), nil
}

func (m *MeasurementsRepository) GetAssetMeasurements(ctx context.Context, assetID string, timeRange measurements.TimeRange) ([]measurements.Measurement, error) {
	ctx, cancel := m.obs.Span(ctx, "measurements.repository.GetAssetMeasurements", zap.String("assetID", assetID))
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.M{"timestamp": -1})

	filter := bson.M{"assetId": assetID}
	timeRangeFilter := bson.M{}

	if timeRange.From != nil {
		timeRangeFilter["$gte"] = timeRange.From
	}

	if timeRange.To != nil {
		timeRangeFilter["$lte"] = timeRange.To
	}

	if len(timeRangeFilter) > 0 {
		filter["timestamp"] = timeRangeFilter
	}

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var dbMeasurements []Measurement
	err = cursor.All(ctx, &dbMeasurements)
	if err != nil {
		return nil, err
	}

	return toMeasurements(dbMeasurements), nil
}

// Aggregation pipeline result
type AveragedMeasurement struct {
	ID            time.Time `bson:"_id"`
	PowerAvg      float64   `bson:"power"`
	StateOfEnergy float64   `bson:"SoE"`
	AssetID       string    `bson:"assetId"`
}

func (m *MeasurementsRepository) GetAssetMeasurementsAveraged(ctx context.Context, assetID string, params measurements.AssetMeasurementAveragedParams) ([]measurements.Measurement, error) {
	ctx, cancel := m.obs.Span(ctx, "measurements.repository.GetAssetMeasurementsAveraged", zap.String("assetID", assetID), zap.Any("params", params))
	defer cancel()

	matchStage := bson.D{
		{"$match", bson.D{
			{"assetId", assetID},
			{"timestamp", bson.D{
				{"$gte", params.From},
				{"$lte", params.To},
			}},
		}},
	}

	// Determine group by instruction
	dateTruncParams, err := mongo2.GroupDateInterval(params.GroupBy)
	if err != nil {
		return nil, err
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{
				{"$dateTrunc", dateTruncParams},
			}},
			{"power", bson.D{{"$avg", "$power.value"}}},
			{"SoE", bson.D{{"$avg", "$stateOfEnergy"}}},
			{"assetId", bson.D{{"$first", "$assetId"}}},
		}},
	}

	val, err := mongo2.SortBy(params.Sort)
	if err != nil {
		return nil, err
	}

	sortStage := bson.D{
		{"$sort",
			bson.D{
				{
					Key: "powerAvg", Value: val,
				},
			},
		},
	}

	pipeline := mongo.Pipeline{matchStage, groupStage, sortStage}
	opts := options.Aggregate()
	cursor, err := m.collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}

	var dbMeasurements []AveragedMeasurement
	err = cursor.All(ctx, &dbMeasurements)
	if err != nil {
		return nil, err
	}

	return toMeasurementsFromAverage(dbMeasurements), nil
}

func toMeasurement(measurement *Measurement) *measurements.Measurement {
	return &measurements.Measurement{
		Time:          measurement.Timestamp,
		Power:         measurement.Power,
		StateOfEnergy: measurement.StateOfEnergy,
	}
}

func toMeasurementsFromAverage(m []AveragedMeasurement) []measurements.Measurement {
	result := make([]measurements.Measurement, len(m))
	for i, measurement := range m {
		result[i] = measurements.Measurement{
			Time:          measurement.ID,
			Power:         measurements.Power{Value: measurement.PowerAvg, Unit: measurements.UnitWatt},
			StateOfEnergy: measurement.StateOfEnergy,
		}
	}
	return result
}

func toMeasurements(m []Measurement) []measurements.Measurement {
	result := make([]measurements.Measurement, len(m))
	for i, measurement := range m {
		result[i] = *toMeasurement(&measurement)
	}
	return result
}

func fromMeasurement(assetId string, measurement *measurements.Measurement) *Measurement {
	return &Measurement{
		AssetID:       assetId,
		Timestamp:     measurement.Time,
		Power:         measurement.Power,
		StateOfEnergy: measurement.StateOfEnergy,
	}
}
