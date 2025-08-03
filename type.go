package statemachine

import "github.com/kkiling/statemachine/internal/storage/sqlite"

type SqliteConfig = sqlite.Config

var NewSqliteStorage = sqlite.NewStorage
