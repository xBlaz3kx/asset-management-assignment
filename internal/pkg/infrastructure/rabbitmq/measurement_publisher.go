package rabbitmq

import (
	"context"
	"encoding/json"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/wagslane/go-rabbitmq"
	"github.com/xBlaz3kx/DevX/observability"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	measurementPublishTopic = "measurements"
	measurementExchange     = "measurements"
)

type MeasurementPublisher struct {
	obs       observability.Observability
	publisher *rabbitmq.Publisher
}

func NewMeasurementPublisher(obs observability.Observability, conn *rabbitmq.Conn) (*MeasurementPublisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		// Enable publisher logging
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsLogger(newLogger(obs)),
		rabbitmq.WithPublisherOptionsExchangeName(measurementExchange),
		rabbitmq.WithPublisherOptionsExchangeDurable,
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	return &MeasurementPublisher{
		obs:       obs.WithSpanKind(trace.SpanKindProducer),
		publisher: publisher,
	}, nil
}

func (p *MeasurementPublisher) Publish(ctx context.Context, measurement measurements.Measurement, assetId string) error {
	ctx, cancel, logger := p.obs.LogSpan(ctx, "measurement.publisher.Publish")
	defer cancel()
	logger.Info("Publishing measurement", zap.Any("measurement", measurement), zap.String("assetId", assetId))

	marshal, err := json.Marshal(measurement)
	if err != nil {
		return err
	}

	headers := rabbitmq.Table{
		"assetId": assetId,
	}

	err = p.publisher.PublishWithContext(
		ctx,
		marshal,
		[]string{measurementPublishTopic},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsHeaders(headers),
	)
	return err
}

func (p *MeasurementPublisher) Close() error {
	p.publisher.Close()
	return nil
}
