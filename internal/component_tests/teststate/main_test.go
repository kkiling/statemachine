package teststate

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kkiling/goplatform/log"
	"github.com/stretchr/testify/require"

	"github.com/kkiling/statemachine/internal/storage/sqlite"
	"github.com/kkiling/statemachine/internal/testutils"
	mock_statemachine "github.com/kkiling/statemachine/mocks"
)

type testDeps struct {
	ctx           context.Context
	ctrl          *gomock.Controller
	service       *StateMachineService
	clock         *mock_statemachine.MockClock
	uuidGenerator *mock_statemachine.MockUUIDGenerator
	storageMock   *mock_statemachine.MockStorage
	storageSqlite *sqlite.Storage
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

func setupTestDB(t *testing.T) *sqlite.Storage {
	// Инициализируем хранилище
	cfg := sqlite.Config{
		DSN: testutils.GetSqliteTestDNS(t), // Берем DSN из переменных окружения
	}
	logger := log.NewLogger(log.DebugLevel)

	s, err := sqlite.NewStorage(cfg, logger)
	require.NoError(t, err)

	// Возвращаем хранилище и функцию очистки
	return s
}

func setupTestDepsSqlite(t *testing.T) *testDeps {
	return setupTestDeps(t, func(deps *testDeps) {
		deps.storageSqlite = setupTestDB(t)
		deps.storageMock = nil
		deps.service = NewState(deps.storageSqlite)
		deps.service.SetClock(deps.clock)
		deps.service.SetUUIDGenerator(deps.uuidGenerator)
	})
}
