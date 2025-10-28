package teststate

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kkiling/goplatform/storagebase/testutils"

	"github.com/kkiling/statemachine/internal/storage/postgresql"
	mock_statemachine "github.com/kkiling/statemachine/mocks"
)

type testDeps struct {
	ctx           context.Context
	ctrl          *gomock.Controller
	service       *StateMachineService
	clock         *mock_statemachine.MockClock
	uuidGenerator *mock_statemachine.MockUUIDGenerator
	storageMock   *mock_statemachine.MockStorage
	storage       *postgresql.Storage
}

func setupTestDeps(t *testing.T, opts ...func(d *testDeps)) *testDeps {
	deps := &testDeps{
		ctx:  context.Background(),
		ctrl: gomock.NewController(t),
	}

	if deps.storageMock == nil {
		deps.storageMock = mock_statemachine.NewMockStorage(deps.ctrl)
	}
	if deps.uuidGenerator == nil {
		deps.uuidGenerator = mock_statemachine.NewMockUUIDGenerator(deps.ctrl)
	}
	if deps.clock == nil {
		deps.clock = mock_statemachine.NewMockClock(deps.ctrl)
	}

	deps.service = NewState(deps.storageMock)
	deps.service.SetClock(deps.clock)
	deps.service.SetUUIDGenerator(deps.uuidGenerator)

	for _, opt := range opts {
		opt(deps)
	}

	return deps
}

func setupTestDepsStorage(t *testing.T) *testDeps {
	return setupTestDeps(t, func(deps *testDeps) {
		deps.storage = postgresql.NewTestStorage(testutils.SetupPostgresqlTestDB(t))
		deps.storageMock = nil
		deps.service = NewState(deps.storage)
		deps.service.SetClock(deps.clock)
		deps.service.SetUUIDGenerator(deps.uuidGenerator)
	})
}
