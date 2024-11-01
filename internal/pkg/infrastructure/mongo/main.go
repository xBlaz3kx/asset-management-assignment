package mongo

import (
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Database = "assets"

func NewClient(connString string) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(connString)
	opts.SetMaxPoolSize(20)
	opts.SetTimeout(time.Second * 10)

	if viper.GetBool("observability.tracing.enabled") {
		// Support for v2 client library is not yet available
		// opts.SetMonitor(otelmongo.NewMonitor())
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
