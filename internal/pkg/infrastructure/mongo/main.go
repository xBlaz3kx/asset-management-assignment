package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Database = "assets"

func NewClient(connString string) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(connString)
	opts.SetMaxPoolSize(20)
	opts.SetTimeout(time.Second * 10)

	// Todo set otel mongo driver wrapper
	// opts.SetMonitor(otelmongo.CommandMonitor{})

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
