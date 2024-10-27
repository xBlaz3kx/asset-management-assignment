package postgres

import (
	"github.com/pkg/errors"
	"github.com/xBlaz3kx/DevX/observability"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func Connect(obs observability.Observability, connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: zapgorm2.New(obs.Log().Logger),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to the database")
	}

	// To consider: If not persisted in the same database, this migration should happen in main.go
	// Migrate the schemas
	err = db.AutoMigrate(&Asset{}, &SimulatorConfiguration{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate schemas")
	}

	return db, nil
}
