package asset_simulation

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

type simulatorManagerTestSuite struct {
	suite.Suite
	manager    *AssetSimulatorManager
	workerMock *MockRunner
}

func (s *simulatorManagerTestSuite) SetupTest() {
	s.manager = NewAssetSimulatorManager(observability.NewNoopObservability())
	s.workerMock = NewMockRunner(s.T())
}

func (s *simulatorManagerTestSuite) TearDownTest() {
	s.workerMock.EXPECT().Stop().Return(nil)
	s.manager.StopAll()
}

func (s *simulatorManagerTestSuite) TestAddWorker() {
	tests := []struct {
		name     string
		workerId string
		err      error
	}{
		{
			name:     "Should add worker",
			workerId: "1",
			err:      nil,
		},
		{
			name:     "Empty worker ID",
			workerId: "",
			err:      ErrWorkerIdCantBeEmpty,
		},
		{
			name:     "Worker already exists",
			workerId: "2",
			err:      ErrWorkerAlreadyExists,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Should add worker":
				s.workerMock.EXPECT().GetId().Return(tt.workerId).Once()
			case "Empty worker ID":
				s.workerMock.EXPECT().GetId().Return(tt.workerId).Once()
			case "Worker already exists":
				s.manager.workers[tt.workerId] = s.workerMock
				s.workerMock.EXPECT().GetId().Return(tt.workerId).Once()
			}

			err := s.manager.AddWorker(s.workerMock)
			if tt.err == nil {
				s.NoError(err)
			} else {
				s.ErrorContains(err, tt.err.Error())
			}
		})
	}
}

func (s *simulatorManagerTestSuite) TestAddAndRunWorker() {
	tests := []struct {
		name     string
		workerId string
		err      error
	}{
		{
			name:     "Should add worker",
			workerId: "1",
			err:      nil,
		},
		{
			name:     "Empty worker ID",
			workerId: "",
			err:      ErrWorkerIdCantBeEmpty,
		},
		{
			name:     "Worker already running",
			workerId: "2",
			err:      nil,
		},
		{
			name:     "Unable to start worker",
			workerId: "3",
			err:      nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Should add worker":
				s.workerMock.EXPECT().GetId().Return(tt.workerId).Once()
				s.workerMock.EXPECT().Start(mock.Anything).Return(nil).Once()
			case "Empty worker ID":
				s.workerMock.EXPECT().GetId().Return(tt.workerId).Once()
			case "Worker already running":
				s.manager.workers[tt.workerId] = s.workerMock
				s.workerMock.EXPECT().IsRunning().Return(true).Once()
				s.workerMock.EXPECT().GetId().Return(tt.workerId).Once()
				s.workerMock.EXPECT().Stop().Return(nil).Once()
				s.workerMock.EXPECT().Start(mock.Anything).Return(nil).Once()
			case "Unable to start worker":
				s.manager.workers[tt.workerId] = s.workerMock
				s.workerMock.EXPECT().IsRunning().Return(false).Once()
				s.workerMock.EXPECT().GetId().Return(tt.workerId).Once()
				s.workerMock.EXPECT().Start(mock.Anything).Return(errors.New("some error")).Once()
			}

			err := s.manager.AddAndStartWorker(context.Background(), s.workerMock)
			if tt.err == nil {
				s.NoError(err)
			} else {
				s.ErrorContains(err, tt.err.Error())
			}
		})
	}
}

func (s *simulatorManagerTestSuite) TestRemoveWorker() {
	tests := []struct {
		name     string
		workerId string
		err      error
	}{
		{
			name:     "Worker found",
			workerId: "1",
		},
		{
			name:     "Worker not found",
			workerId: "2",
			err:      ErrWorkerDoesntExist,
		},
		{
			name:     "Worker not running",
			workerId: "3",
			err:      nil,
		},
		{
			name:     "Unable to stop worker",
			workerId: "4",
			err:      errors.New("unable to stop worker"),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Worker found":
				s.manager.workers[tt.workerId] = s.workerMock
				s.workerMock.EXPECT().IsRunning().Return(true).Once()
				s.workerMock.EXPECT().Stop().Return(nil).Once()
			case "Worker not found":
			case "Worker not running":
				s.manager.workers[tt.workerId] = s.workerMock
				s.workerMock.EXPECT().IsRunning().Return(false).Once()
			case "Unable to stop worker":
				s.manager.workers[tt.workerId] = s.workerMock
				s.workerMock.EXPECT().IsRunning().Return(true).Once()
				s.workerMock.EXPECT().Stop().Return(errors.New("some error")).Once()
			}

			err := s.manager.RemoveWorker(tt.workerId)
			if tt.err == nil {
				s.NoError(err)
			} else {
				s.ErrorContains(err, tt.err.Error())
			}
		})
	}
}

func (s *simulatorManagerTestSuite) TestGetWorker() {
	s.T().Skip()
	tests := []struct {
		name     string
		workerId string
		isFound  bool
	}{
		{
			name:     "Worker found",
			workerId: "1",
		},
		{
			name:     "Worker not found",
			workerId: "2",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Worker found":
				s.manager.workers[tt.workerId] = s.workerMock
			}

			worker, isFound := s.manager.GetWorker(tt.workerId)
			if tt.isFound {
				s.Equal(s.workerMock, worker)
				s.True(isFound)
			} else {
				s.Nil(worker)
				s.False(isFound)
			}
		})
	}
}

func (s *simulatorManagerTestSuite) TestStartWorkers() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *simulatorManagerTestSuite) TestGetWorkers() {
	s.T().Skip()

	tests := []struct {
		name string
	}{
		{
			name: "One worker",
		},
		{
			name: "No workers",
		},
		{
			name: "Multiple workers",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			expectedWorkers := []Runner{}
			switch tt.name {
			case "One worker":
				s.manager.workers["workerId"] = s.workerMock
				expectedWorkers = append(expectedWorkers, s.workerMock)
			case "Multiple workers":
				mock1 := NewMockRunner(s.T())
				s.manager.workers["workerId1"] = mock1
				expectedWorkers = append(expectedWorkers, mock1)

				mock2 := NewMockRunner(s.T())
				s.manager.workers["workerId2"] = mock2
				expectedWorkers = append(expectedWorkers, mock1)

				mock3 := NewMockRunner(s.T())
				s.manager.workers["workerId3"] = mock3
				expectedWorkers = append(expectedWorkers, mock1)

			}

			workers := s.manager.GetWorkers()
			s.Len(workers, len(expectedWorkers))
		})
	}
}

func (s *simulatorManagerTestSuite) TestStopAll() {
	tests := []struct {
		name string
	}{
		{
			name: "One worker",
		},
		{
			name: "No workers",
		},
		{
			name: "Multiple workers",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "One worker":
				s.manager.workers["workerId"] = s.workerMock
				s.workerMock.EXPECT().Stop().Return(nil)
			case "Multiple workers":
				mock1 := NewMockRunner(s.T())
				mock1.EXPECT().Stop().Return(nil)

				mock2 := NewMockRunner(s.T())
				mock2.EXPECT().Stop().Return(nil)

				mock3 := NewMockRunner(s.T())
				mock3.EXPECT().Stop().Return(nil)

				s.manager.workers["workerId1"] = mock1
				s.manager.workers["workerId2"] = mock2
				s.manager.workers["workerId3"] = mock3
			}

			s.manager.StopAll()
			if tt.name == "No workers" {
				s.workerMock.AssertCalled(t, "Stop")
			}
		})
	}
}

func TestSimulatorManager(t *testing.T) {
	suite.Run(t, new(simulatorManagerTestSuite))
}
