package statemachine

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kkiling/goplatform/storagebase"
)

type Config struct {
}

type StateMachine[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string, CreateOptionsT CreateOptions] struct {
	cfg           Config
	runner        Runner[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]
	storage       Storage
	clock         Clock
	uuidGenerator UUIDGenerator
}

func NewService[DataT any, FailDataT any, MetaDataT any, StepT ~string, TypeT ~string, CreateOptionsT CreateOptions](
	cfg Config,
	storage Storage,
	runner Runner[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT],
) *StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT] {

	sm := StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]{
		cfg:           cfg,
		runner:        runner,
		storage:       storage,
		clock:         &realClock{},
		uuidGenerator: &uuidGenerator{},
	}

	return &sm
}

func (i *StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]) getStateByIdempotencyKey(
	ctx context.Context,
	idempotencyKey string,
) (*State[DataT, FailDataT, MetaDataT, StepT, TypeT], error) {
	findState, err := i.storage.GetStateByIdempotencyKey(ctx, idempotencyKey)
	switch {
	case err == nil:
	case errors.Is(err, storagebase.ErrNotFound): // Стейт не найден
		return nil, nil
	default:
		return nil, fmt.Errorf("storage.GetStateByIdempotencyKey: %w", err)
	}
	// Стейт найден
	return mapStorageToState[DataT, FailDataT, MetaDataT, StepT, TypeT](findState)
}

func (i *StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]) GetStateByID(
	ctx context.Context,
	stateID uuid.UUID,
) (*State[DataT, FailDataT, MetaDataT, StepT, TypeT], error) {
	findState, err := i.storage.GetStateByID(ctx, stateID)
	switch {
	case err == nil:
	case errors.Is(err, storagebase.ErrNotFound): // Стейт не найден
		return nil, nil
	default:
		return nil, fmt.Errorf("storage.GetStateByID: %w", err)
	}
	return mapStorageToState[DataT, FailDataT, MetaDataT, StepT, TypeT](findState)
}

// Create создание стейт машины
func (i *StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]) Create(
	ctx context.Context,
	options CreateOptionsT,
) (*State[DataT, FailDataT, MetaDataT, StepT, TypeT], error) {
	// Проверяем стейт на наличие ключа идемпотентности
	if findState, err := i.getStateByIdempotencyKey(ctx, options.GetIdempotencyKey()); err != nil {
		return nil, fmt.Errorf("getStateByIdempotencyKey: %w", err)
	} else if findState != nil {
		return findState, ErrAlreadyExists
	}

	now := i.clock.Now()
	// Вызываем раннер, который выполняет бизнес логику и возвращает issueData
	create, err := i.runner.Create(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("deliverystate.Create: %w", err)
	}

	newIssue := State[DataT, FailDataT, MetaDataT, StepT, TypeT]{
		ID:             i.uuidGenerator.New(),
		IdempotencyKey: options.GetIdempotencyKey(),
		CreatedAt:      now,
		UpdatedAt:      now,
		Status:         NewStatus,
		Step:           create.FirstStep,
		Type:           i.runner.Type(),
		Data:           create.Data,
		MetaData:       create.MetaData,
	}

	newStorageState, err := mapStateToStorage[DataT, FailDataT, MetaDataT, StepT, TypeT](&newIssue)
	if err != nil {
		return nil, fmt.Errorf("mapStateToStorage: %w", err)
	}

	saveErr := i.storage.CreateState(ctx, newStorageState)
	if saveErr != nil {
		return nil, fmt.Errorf("storage.CreateState: %w", saveErr)
	}

	return &newIssue, nil
}

func (i *StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]) initStepper() *Stepper[DataT, FailDataT, MetaDataT, StepT, TypeT] {
	stepper := NewStepper[DataT, FailDataT, MetaDataT, StepT, TypeT](i.storage, i.clock)
	stepsRegistration := i.runner.StepRegistration(StepRegistrationParams{})
	for s, step := range stepsRegistration.Steps {
		stepper.Add(s, step)
	}
	return stepper
}

// Complete выполнение стейта
func (i *StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]) Complete(
	ctx context.Context,
	stateID uuid.UUID,
	options ...any,
) (st *State[DataT, FailDataT, MetaDataT, StepT, TypeT], executeErr error, err error) {
	// Проверяем стейт на наличие ключа идемпотентности
	findState, err := i.GetStateByID(ctx, stateID)
	if err != nil {
		return nil, nil, fmt.Errorf("getStateByID: %w", err)
	}

	if findState == nil {
		return nil, nil, fmt.Errorf("state not found: %w", ErrNotFound)
	}

	if findState.Status == FailedStatus || findState.Status == CompletedStatus {
		return nil, nil, ErrInTerminalStatus
	}

	stepper := i.initStepper()
	res, eErr, err := stepper.Compete(ctx, *findState, options...)
	if err != nil {
		return nil, nil, fmt.Errorf("stepper.Compete: %w", err)
	}

	return res, eErr, nil
}

// SetClock устанавливает кастомную реализацию часов
func (i *StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]) SetClock(clock Clock) {
	i.clock = clock
}

// SetUUIDGenerator устанавливает кастомную реализацию uuid генератора
func (i *StateMachine[DataT, FailDataT, MetaDataT, StepT, TypeT, CreateOptionsT]) SetUUIDGenerator(generator UUIDGenerator) {
	i.uuidGenerator = generator
}
