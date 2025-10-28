package postgresql

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/statemachine/internal/storage"
	"github.com/kkiling/statemachine/internal/storage/statemachine"
)

func (s *Storage) CreateState(ctx context.Context, state *storage.State) error {
	queries := s.getQueries(ctx)

	params := statemachine.CreateStateParams{
		ID:             state.ID,
		IdempotencyKey: state.IdempotencyKey,
		CreatedAt:      state.CreatedAt,
		UpdatedAt:      state.UpdatedAt,
		Status:         int(state.Status),
		Step:           state.Step,
		Type:           state.Type,
		Data:           state.Data,
		FailData: func() []byte {
			if len(state.FailData) == 0 {
				return nil
			}
			return state.FailData
		}(),
		MetaData: func() []byte {
			if len(state.MetaData) == 0 {
				return nil
			}
			return state.MetaData
		}(),
	}
	err := queries.CreateState(ctx, params)

	return s.base.HandleError(err)
}

func (s *Storage) GetStateByIdempotencyKey(ctx context.Context, idempotencyKey string) (*storage.State, error) {
	queries := s.getQueries(ctx)

	res, err := queries.GetStateByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	return &storage.State{
		ID:             res.ID,
		IdempotencyKey: res.IdempotencyKey,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
		Status:         uint8(res.Status),
		Step:           res.Step,
		Type:           res.Type,
		Data:           res.Data,
		FailData:       res.FailData,
		MetaData:       res.MetaData,
		Error:          res.Error,
	}, nil
}

func (s *Storage) GetStateByID(ctx context.Context, stateID uuid.UUID) (*storage.State, error) {
	queries := s.getQueries(ctx)

	res, err := queries.GetStateByID(ctx, stateID)
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	return &storage.State{
		ID:             res.ID,
		IdempotencyKey: res.IdempotencyKey,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
		Status:         uint8(res.Status),
		Step:           res.Step,
		Type:           res.Type,
		Data:           res.Data,
		FailData:       res.FailData,
		MetaData:       res.MetaData,
		Error:          res.Error,
	}, nil
}

func (s *Storage) SaveStepExecuteInfo(ctx context.Context, execute storage.StepExecuteInfo) error {
	queries := s.getQueries(ctx)

	err := queries.SaveStepExecuteInfo(ctx, statemachine.SaveStepExecuteInfoParams{
		StateID:            execute.StateID,
		StartExecutedAt:    execute.StartExecutedAt,
		CompleteExecutedAt: execute.CompleteExecutedAt,
		Error:              execute.Error,
		PreviewStep:        execute.PreviewStep,
		NextStep:           execute.NextStep,
	})

	return s.base.HandleError(err)
}

func (s *Storage) GetStepExecuteInfos(ctx context.Context, stateID uuid.UUID) ([]storage.StepExecuteInfo, error) {
	queries := s.getQueries(ctx)
	res, err := queries.GetStepExecuteInfos(ctx, stateID)
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	return lo.Map(res, func(item statemachine.GetStepExecuteInfosRow, index int) storage.StepExecuteInfo {
		return storage.StepExecuteInfo{
			StateID:            item.StateID,
			StartExecutedAt:    item.StartExecutedAt,
			CompleteExecutedAt: item.CompleteExecutedAt,
			Error:              item.Error,
			PreviewStep:        item.PreviewStep,
			NextStep:           item.NextStep,
		}
	}), nil
}

func (s *Storage) UpdateState(ctx context.Context, stateID uuid.UUID, state storage.UpdateState) error {
	queries := s.getQueries(ctx)

	_, err := queries.UpdateState(ctx, statemachine.UpdateStateParams{
		UpdatedAt: state.UpdatedAt,
		Status:    int(state.Status),
		Step:      state.Step,
		Data:      state.Data,
		FailData: func() []byte {
			if len(state.FailData) == 0 {
				return nil
			}
			return state.FailData
		}(),
		MetaData: func() []byte {
			if len(state.MetaData) == 0 {
				return nil
			}
			return state.MetaData
		}(),
		Error: state.Error,
		ID:    stateID,
	})

	return s.base.HandleError(err)
}
