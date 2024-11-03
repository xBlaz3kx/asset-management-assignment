package mongo

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GroupDateInterval(groupBy string) (bson.D, error) {
	// Determine groupby instruction
	dateTruncParams := bson.D{
		{"date", "$timestamp"},
	}
	switch groupBy {
	case "minute":
		dateTruncParams = append(dateTruncParams, bson.E{Key: "unit", Value: "minute"})
	case "15min":
		dateTruncParams = append(dateTruncParams, bson.E{Key: "unit", Value: "minute"}, bson.E{Key: "minute", Value: 15})
	case "hour":
		dateTruncParams = append(dateTruncParams, bson.E{Key: "unit", Value: "hour"})
	default:
		return nil, errors.New("invalid groupBy value")
	}

	return dateTruncParams, nil
}

func SortBy(sort string) (int, error) {
	switch sort {
	case "asc":
		return 1, nil
	case "desc":
		return -1, nil
	default:
		return 0, errors.New("invalid sort value")
	}
}
