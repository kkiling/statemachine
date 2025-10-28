package statemachine

import (
	"github.com/kkiling/goplatform/storagebase/postgrebase"
)

type StorageConfig = postgrebase.Config

var NewStorage = postgrebase.NewStorage

var NewPgConn = postgrebase.NewPgConn
