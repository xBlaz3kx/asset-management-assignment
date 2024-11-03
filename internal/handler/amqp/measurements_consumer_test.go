package amqp

import (
	"context"
	"testing"
	"time"

	serviceMock "asset-measurements-assignment/internal/domain/measurements/service/mocks"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wagslane/go-rabbitmq"
	"github.com/xBlaz3kx/DevX/observability"
)

func TestHandler_handleMeasurement(t *testing.T) {
	mockObs := observability.NewNoopObservability()

	tests := []struct {
		name   string
		args   rabbitmq.Delivery
		result rabbitmq.Action
	}{
		{
			name: "Valid measurement",
			args: rabbitmq.Delivery{
				Delivery: amqp091.Delivery{
					Headers: amqp091.Table{
						"assetId": "1",
					},
					ContentType: "application/json",
					MessageId:   uuid.New().String(),
					Timestamp:   time.Now(),
					Exchange:    measurementExchange,
					RoutingKey:  measurementRoutingKey,
					Body:        []byte(`{"power": {"value": 1000, "unit": "W"}, "time": "2021-09-01T12:00:00Z", "stateOfEnergy": 1.00}`),
				},
			},
			result: rabbitmq.Ack,
		},
		{
			name: "AssetId is missing",
			args: rabbitmq.Delivery{
				Delivery: amqp091.Delivery{
					Headers:     amqp091.Table{},
					ContentType: "application/json",
					MessageId:   uuid.New().String(),
					Timestamp:   time.Now(),
					Exchange:    measurementExchange,
					RoutingKey:  measurementRoutingKey,
					Body:        []byte(`{"power": {"value": 1000.0, "unit": "W"}, "time": "2021-09-01T12:00:00Z", "stateOfEnergy": 1.00}`),
				},
			},
			result: rabbitmq.NackDiscard,
		},
		{
			name: "Body is not a valid JSON",
			args: rabbitmq.Delivery{
				Delivery: amqp091.Delivery{
					Headers: amqp091.Table{
						"assetId": "1",
					},
					ContentType: "application/json",
					MessageId:   uuid.New().String(),
					Timestamp:   time.Now(),
					Exchange:    measurementExchange,
					RoutingKey:  measurementRoutingKey,
					Body:        []byte(`{"power": {"value": 1000, "unit": "W"}, "time": "2021-09-01T12:00:00Z", "stateOfEnergy": 1.00`),
				},
			},
			result: rabbitmq.NackDiscard,
		},
		{
			name: "Unable to store measurement",
			args: rabbitmq.Delivery{
				Delivery: amqp091.Delivery{
					Headers: amqp091.Table{
						"assetId": "3",
					},
					ContentType: "application/json",
					MessageId:   uuid.New().String(),
					Timestamp:   time.Now(),
					Exchange:    measurementExchange,
					RoutingKey:  measurementRoutingKey,
					Body:        []byte(`{"power": {"value": 1000, "unit": "W"}, "time": "2021-09-01T12:00:00Z", "stateOfEnergy": 1.00}`),
				},
			},
			result: rabbitmq.NackDiscard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumerServiceMock := serviceMock.NewMockConsumerService(t)

			switch tt.name {
			case "Valid measurement":
				consumerServiceMock.EXPECT().AddMeasurement(mock.Anything, "1", mock.Anything).Return(nil).Once()
			case "Unable to store measurement":
				consumerServiceMock.EXPECT().
					AddMeasurement(mock.Anything, "3", mock.Anything).
					Return(errors.New("failed to store measurement")).Once()
			}

			h := &Handler{
				obs:     mockObs,
				service: consumerServiceMock,
			}

			handleFn := h.handleMeasurement(context.Background())
			result := handleFn(tt.args)

			assert.Equal(t, tt.result, result)
		})
	}
}
