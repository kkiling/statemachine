package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkiling/goplatform/storagebase/postgrebase"

	"github.com/kkiling/statemachine/internal/storage/statemachine"
)

type Storage struct {
	base *postgrebase.Storage
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		base: postgrebase.NewStorage(pool),
	}
}

func (s *Storage) getQueries(ctx context.Context) *statemachine.Queries {
	return statemachine.New(s.base.Next(ctx))
}

func (s *Storage) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	return s.base.RunTransaction(ctx, txFunc)
}

func NewTestStorage(base *postgrebase.Storage) *Storage {
	return &Storage{
		base: base,
	}
}
