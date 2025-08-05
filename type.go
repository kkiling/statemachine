package statemachine

import (
	"github.com/kkiling/goplatform/storagebase/sqlitebase"

	"github.com/kkiling/statemachine/internal/storage/sqlite"
)

type SqliteConfig = sqlitebase.Config

var NewSqliteStorage = sqlite.NewStorage
