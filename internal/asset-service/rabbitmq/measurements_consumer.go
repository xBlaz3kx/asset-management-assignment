package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
	"asset-measurements-assignment/internal/domain/measurements/service"
	rmq "asset-measurements-assignment/internal/pkg/infrastructure/rabbitmq"
	"github.com/wagslane/go-rabbitmq"
	"github.com/xBlaz3kx/DevX/observability"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const measurementRoutingKey = "measurement"
const measurementExchange = "measurement"

type Handler struct {
	obs      observability.Observability
	consumer *rabbitmq.Consumer
	service  service.ConsumerService
}

func NewHandler(obs observability.Observability, conn *rabbitmq.Conn, service service.ConsumerService) (*Handler, error) {
	// Create a new measurements consumer
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"",
		// Enable consumer logging
		rabbitmq.WithConsumerOptionsLogger(rmq.NewLogger(obs)),
		rabbitmq.WithConsumerOptionsRoutingKey(measurementRoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(measurementExchange),
		rabbitmq.WithConsumerOptionsQueueDurable,
	)
	if err != nil {
		return nil, err
	}

	return &Handler{
		service:  service,
		obs:      obs.WithSpanKind(trace.SpanKindConsumer),
		consumer: consumer,
	}, nil
}

func (h *Handler) Start(ctx context.Context) error {
	// Start consuming messages in a separate goroutine
	go func() {
		err := h.consumer.Run(h.handleMeasurement(ctx))
		if err != nil {
			h.obs.Log().Error("Consumer failed", zap.Error(err))
		}
	}()

	return nil
}

// handleMeasurement handles the incoming measurement messages.
// It unmarshal the message, gets the assetId from the header and attempts to store the measurement.
func (h *Handler) handleMeasurement(ctx context.Context) func(d rabbitmq.Delivery) (action rabbitmq.Action) {
	return func(delivery rabbitmq.Delivery) (action rabbitmq.Action) {
		consumeCtx, cancel, logger := h.obs.LogSpanWithTimeout(ctx, "measurement.consumer.Handle",
			time.Second*10,
		)
		defer cancel()
		logger.Info("Consuming measurement")

		// Check if the content type is JSON
		if delivery.ContentType != "application/json" {
			return rabbitmq.NackDiscard
		}

		// Unmarshal the message
		var measurement measurements.Measurement
		err := json.Unmarshal(delivery.Body, &measurement)
		if err != nil {
			return rabbitmq.NackDiscard
		}

		// Get assetId from header
		assetID, ok := delivery.Headers["assetId"].(string)
		if !ok || assetID == "" {
			return rabbitmq.NackDiscard
		}

		// Store the measurement
		err = h.service.AddMeasurement(consumeCtx, assetID, measurement)
		if err != nil {
			logger.With(zap.Error(err)).Error("Failed to store measurement")
			// Requeue?
			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	}
}

// Close closes the consumer.
func (h *Handler) Close() error {
	h.consumer.Close()
	return nil
}
