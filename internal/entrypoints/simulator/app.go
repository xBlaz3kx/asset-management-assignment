package simulator

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain/simulator"
	"asset-measurements-assignment/internal/pkg/asset_simulation"
	"asset-measurements-assignment/internal/pkg/infrastructure/postgres"
	"asset-measurements-assignment/internal/pkg/infrastructure/rabbitmq"
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
	PostgresConnection string `yaml:"postgres" mapstructure:"postgres"`

	// RabbitMQ connection string
	// Example: amqp://guest:guest@localhost:5672/
	RabbitMQConnection string `yaml:"rabbitmq" mapstructure:"rabbitmq"`

	// Observability settings
	Observability observability.Config `yaml:"observability" mapstructure:"observability"`
}

var serviceInfo = observability.ServiceInfo{
	Name:    "simulator",
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

	obs.Log().Info("Starting simulator", zap.Any("config", cfg))

	env := viper.GetString("environment")
	if env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to Postgres
	postgresDb, err := postgres.Connect(obs, cfg.PostgresConnection)
	if err != nil {
		return err
	}

	rabbitmqConn, err := goRabbit.NewConn(cfg.RabbitMQConnection,
		goRabbit.WithConnectionOptionsLogging,
		goRabbit.WithConnectionOptionsReconnectInterval(time.Second*3),
	)
	if err != nil {
		return err
	}
	defer rabbitmqConn.Close()

	// Create new simulator configuration repository
	configRepository := postgres.NewSimulatorConfigurationRepository(obs, postgresDb)

	// Fetch assets from Postgres
	configs, err := configRepository.GetConfigurations(ctx)
	if err != nil {
		return err
	}

	// Create measurements publisher
	measurementPublisher, err := rabbitmq.NewMeasurementPublisher(obs, rabbitmqConn)
	if err != nil {
		return err
	}

	// Create new asset simulator worker manager
	workerManager := asset_simulation.NewAssetSimulatorManager(obs)

	// Add workers based on fetched configurations
	createWorkersFromConfigurations(obs, workerManager, configs, measurementPublisher)

	// Start asset simulator workers
	workerManager.StartWorkers(ctx)

	// Create HTTP server for healthchecks
	httpServer := devxHttp.NewServer(devxHttp.Configuration{Address: ":80"}, obs)
	go func() {
		httpServer.Run()
	}()

	<-ctx.Done()
	obs.Log().Info("Shutting down simulator")
	workerManager.StopAll()

	return nil
}

func createWorkersFromConfigurations(
	obs observability.Observability,
	workerManager *asset_simulation.AssetSimulatorManager,
	configs []simulator.Configuration,
	measurementPublisher *rabbitmq.MeasurementPublisher,
) {
	for _, config := range configs {

		configuration := asset_simulation.Configuration{
			AssetId:             config.AssetId,
			MinPower:            config.MinPower,
			MaxPower:            config.MaxPower,
			MaxPowerStep:        config.MaxPowerStep,
			MeasurementInterval: config.MeasurementInterval,
		}

		worker, err := asset_simulation.NewSimpleAssetSimulator(obs, configuration, measurementPublisher)
		if err != nil {
			obs.Log().Error("Failed to create worker", zap.Error(err))
			continue
		}

		workerManager.AddWorker(worker)
	}
}
