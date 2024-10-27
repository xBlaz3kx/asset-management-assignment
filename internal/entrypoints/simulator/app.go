package simulator

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/pkg/infrastructure/postgres"
	"asset-measurements-assignment/internal/pkg/infrastructure/rabbitmq"
	"asset-measurements-assignment/internal/pkg/simulator_worker"
	"github.com/GLCharge/otelzap"
	goRabbit "github.com/wagslane/go-rabbitmq"
	devxHttp "github.com/xBlaz3kx/DevX/http"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type Config struct {
	// Postgres connection string
	PostgresConnection string `yaml:"postgres" mapstructure:"postgres"`

	// RabbitMQ connection string
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
	workerManager := simulator_worker.NewAssetSimulatorManager()

	// Add workers based on fetched configurations
	for _, config := range configs {
		worker := simulator_worker.NewSimpleAssetSimulator(obs, config, measurementPublisher)
		workerManager.AddWorker(worker)
	}

	// Start asset simulator workers
	workerManager.StartWorkers(ctx)

	// Create HTTP server for healthchecks
	httpServer := devxHttp.NewServer(devxHttp.Configuration{Address: ":80"}, obs)
	go func() {
		httpServer.Run()
	}()

	<-ctx.Done()
	workerManager.StopAll()

	return nil
}
