package postgres

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/DevX/observability"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"moul.io/zapgorm2"
)

func Connect(obs observability.Observability, connectionString string) (*gorm.DB, error) {
	logger := zapgorm2.New(obs.Log().Logger)
	logger.LogLevel = gormlogger.Info
	logger.SetAsDefault()

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to the database")
	}

	// Check if tracing is enabled from env
	if viper.GetBool("observability.tracing.enabled") {
		if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
			return nil, errors.Wrap(err, "failed to use tracing plugin")
		}
	}

	// To consider: If not persisted in the same database, this migration should happen in main.go
	// Migrate the schemas
	err = db.AutoMigrate(&Asset{}, &SimulatorConfiguration{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate schemas")
	}

	return db, nil
}
