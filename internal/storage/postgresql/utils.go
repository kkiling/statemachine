package postgresql

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func toTimePtr(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func toTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t != nil {
		return pgtype.Timestamptz{
			Time:  *t,
			Valid: true,
		}
	}
	return pgtype.Timestamptz{}
}
