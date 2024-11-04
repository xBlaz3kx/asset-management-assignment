package simulator

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain/simulator/service"
	"asset-measurements-assignment/internal/pkg/infrastructure/postgres"
	assetSimulation "asset-measurements-assignment/internal/simulator/asset_simulation"
	"asset-measurements-assignment/internal/simulator/http"
	postgres2 "asset-measurements-assignment/internal/simulator/postgres"
	"asset-measurements-assignment/internal/simulator/rabbitmq"
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
	configRepository := postgres2.NewSimulatorConfigurationRepository(obs, postgresDb)

	// Create new asset simulator worker manager
	workerManager := assetSimulation.NewAssetSimulatorManager(obs)

	// Create measurements publisher
	measurementPublisher, err := rabbitmq.NewMeasurementPublisher(obs, rabbitmqConn)
	if err != nil {
		return err
	}

	// Create new asset configuration service
	configService := service.NewConfigService(obs, configRepository, workerManager, measurementPublisher)
	err = configService.StartWorkersFromDatabaseConfigurations(ctx)
	if err != nil {
		// Log error and continue
		obs.Log().With(zap.Error(err)).Error("Failed to start workers from database configurations")
	}

	// Create HTTP server for healthchecks
	httpServer := devxHttp.NewServer(devxHttp.Configuration{Address: ":80"}, obs)
	router := httpServer.Router()
	configHandler := http.NewSimulatorConfigHandler(configService)
	configHandler.RegisterRoutes(router)

	go func() {
		httpServer.Run()
	}()

	<-ctx.Done()
	obs.Log().Info("Shutting down simulator")
	workerManager.StopAll()

	return nil
}
