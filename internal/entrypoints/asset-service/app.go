package asset_service

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain/assets"
	measurements "asset-measurements-assignment/internal/domain/measurements/service"
	"asset-measurements-assignment/internal/handler/amqp"
	"asset-measurements-assignment/internal/handler/http"
	"asset-measurements-assignment/internal/pkg/infrastructure/mongo"
	"asset-measurements-assignment/internal/pkg/infrastructure/postgres"
	"github.com/GLCharge/otelzap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	goRabbit "github.com/wagslane/go-rabbitmq"
	devxHttp "github.com/xBlaz3kx/DevX/http"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type Config struct {
	// Postgres connection string
	// Example: postgres://user:password@localhost:5432/dbname?sslmode=disable
	Postgres string `yaml:"postgres" mapstructure:"postgres" json:"postgres"`

	// RabbitMQ connection string
	// Example: amqp://guest:guest@localhost:5672/
	Rabbitmq string `yaml:"rabbitmq" mapstructure:"rabbitmq" json:"rabbitmq"`

	// MongoDB connection string
	// Example: mongodb://localhost:27017
	Mongo string `yaml:"mongo" mapstructure:"mongo" json:"mongo"`

	// HTTP server settings
	Http devxHttp.Configuration `yaml:"http" mapstructure:"http" json:"http"`

	// Observability settings
	Observability observability.Config `yaml:"observability" mapstructure:"observability" json:"observability"`
}

var serviceInfo = observability.ServiceInfo{
	Name:    "asset-service",
	Version: "0.0.1",
}

func Run(ctx context.Context, cfg Config) error {
	// Initialize observability
	obs, err := observability.NewObservability(ctx, serviceInfo, cfg.Observability)
	if err != nil {
		otelzap.L().With(zap.Error(err)).Fatal("Could not initialize observability")
	}
	defer func() {
		shutdownCtx, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel2()
		_ = obs.Shutdown(shutdownCtx)
	}()

	obs.Log().Info("Starting asset service", zap.Any("config", cfg))

	env := viper.GetString("environment")
	if env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to Postgres
	postgresDb, err := postgres.Connect(obs, cfg.Postgres)
	if err != nil {
		return err
	}

	// Create asset repository
	assetRepository := postgres.NewAssetRepository(obs, postgresDb)

	// Connect to MongoDB
	mongoClient, err := mongo.NewClient(cfg.Mongo)
	if err != nil {
		return err
	}
	defer func() {
		_ = mongoClient.Disconnect(ctx)
	}()
	mongoDb := mongoClient.Database(mongo.Database)

	measurementsRepository, err := mongo.NewMeasurementsRepository(obs, mongoDb)
	if err != nil {
		return err
	}

	// Create consumer service
	consumerService := measurements.NewConsumerService(obs, measurementsRepository, assetRepository)

	rabbitMqConn, err := goRabbit.NewConn(cfg.Rabbitmq,
		goRabbit.WithConnectionOptionsLogging,
		goRabbit.WithConnectionOptionsReconnectInterval(time.Second*5),
	)
	if err != nil {
		return err
	}
	defer rabbitMqConn.Close()

	// Create rabbitmq consumer
	consumer, err := amqp.NewHandler(obs, rabbitMqConn, consumerService)
	if err != nil {
		return err
	}
	defer consumer.Close()

	err = consumer.Start(ctx)
	if err != nil {
		return err
	}

	// Create asset service
	assetService := assets.NewService(obs, assetRepository)

	// Create measurements service
	measurementsService := measurements.NewMeasurementsService(obs, measurementsRepository)

	// Create HTTP server
	httpServer := devxHttp.NewServer(cfg.Http, obs)
	router := httpServer.Router()

	// Asset handler
	handler := http.NewAssetGinHandler(assetService)
	handler.RegisterRoutes(router)

	// Measurements handler
	measurementsGinHandler := http.NewMeasurementsGinHandler(measurementsService)
	measurementsGinHandler.RegisterRoutes(router)

	go func() {
		// Todo define health checks
		httpServer.Run()
	}()

	<-ctx.Done()
	obs.Log().Info("Shutting down asset service")
	httpServer.Shutdown()

	return nil
}
