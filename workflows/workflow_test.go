package workflows_test

import (
	"eth-temporal/app"
	"eth-temporal/app/activities"
	"eth-temporal/app/workflows"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_GetLatestBlockNumWorkflow() {
	s.env.OnActivity(activities.GetLatestBlockNum, mock.Anything).Return(uint64(1024), nil)
	s.env.ExecuteWorkflow(workflows.GetLatestBlockNumWorkflow)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	var result uint64
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(uint64(1024), result)
}

func (s *UnitTestSuite) Test_GetBlockWorkflow() {
	s.env.OnActivity(activities.GetBlockByNumber, mock.Anything, uint64(1024)).Return(app.Block{Number: uint64(1024)}, nil)
	s.env.OnActivity(activities.UpsertBlockToPostgres, mock.Anything, mock.Anything).Return(nil)
	s.env.ExecuteWorkflow(workflows.GetBlockWorkflow, uint64(1024))

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var block app.Block
	s.NoError(s.env.GetWorkflowResult(&block))
	s.Equal(uint64(1024), block.Number)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
