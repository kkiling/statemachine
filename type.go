package statemachine

import (
	"github.com/kkiling/goplatform/storagebase/postgrebase"
	"github.com/kkiling/statemachine/internal/storage/postgresql"
)

type StorageConfig = postgrebase.Config

var NewPgConn = postgrebase.NewPgConn

var NewStorage = postgresql.NewStorage
