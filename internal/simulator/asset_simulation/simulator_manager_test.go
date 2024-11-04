package asset_simulation

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

type simulatorManagerTestSuite struct {
	suite.Suite
	manager *AssetSimulatorManager
}

func (s *simulatorManagerTestSuite) SetupTest() {
	s.manager = NewAssetSimulatorManager(observability.NewNoopObservability())
}

func (s *simulatorManagerTestSuite) TearDownTest() {
	s.manager.StopAll()
}

func (s *simulatorManagerTestSuite) TestAddWorker() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *simulatorManagerTestSuite) TestRemoveWorker() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *simulatorManagerTestSuite) TestGetWorker() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

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
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestSimulatorManager(t *testing.T) {
	t.Skip("Not implemented")
	suite.Run(t, new(simulatorManagerTestSuite))
}
