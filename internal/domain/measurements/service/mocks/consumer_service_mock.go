// Code generated by mockery v2.46.3. DO NOT EDIT.

package service

import (
	measurements "asset-measurements-assignment/internal/domain/measurements"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockConsumerService is an autogenerated mock type for the ConsumerService type
type MockConsumerService struct {
	mock.Mock
}

type MockConsumerService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConsumerService) EXPECT() *MockConsumerService_Expecter {
	return &MockConsumerService_Expecter{mock: &_m.Mock}
}

// AddMeasurement provides a mock function with given fields: ctx, assetId, measurement
func (_m *MockConsumerService) AddMeasurement(ctx context.Context, assetId string, measurement measurements.Measurement) error {
	ret := _m.Called(ctx, assetId, measurement)

	if len(ret) == 0 {
		panic("no return value specified for AddMeasurement")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, measurements.Measurement) error); ok {
		r0 = rf(ctx, assetId, measurement)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockConsumerService_AddMeasurement_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddMeasurement'
type MockConsumerService_AddMeasurement_Call struct {
	*mock.Call
}

// AddMeasurement is a helper method to define mock.On call
//   - ctx context.Context
//   - assetId string
//   - measurement measurements.Measurement
func (_e *MockConsumerService_Expecter) AddMeasurement(ctx interface{}, assetId interface{}, measurement interface{}) *MockConsumerService_AddMeasurement_Call {
	return &MockConsumerService_AddMeasurement_Call{Call: _e.mock.On("AddMeasurement", ctx, assetId, measurement)}
}

func (_c *MockConsumerService_AddMeasurement_Call) Run(run func(ctx context.Context, assetId string, measurement measurements.Measurement)) *MockConsumerService_AddMeasurement_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(measurements.Measurement))
	})
	return _c
}

func (_c *MockConsumerService_AddMeasurement_Call) Return(_a0 error) *MockConsumerService_AddMeasurement_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConsumerService_AddMeasurement_Call) RunAndReturn(run func(context.Context, string, measurements.Measurement) error) *MockConsumerService_AddMeasurement_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockConsumerService creates a new instance of MockConsumerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConsumerService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConsumerService {
	mock := &MockConsumerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}