package query

import "EasyGo/config"

const (
	INSERT = iota
	DELETE
	UPDATE
	FETCH
	FETCHRAW
)

type Processor interface {
	Exec(scene int, cmd string, param ...interface{}) (interface{}, error)
	Connect(conf config.Result) error
	Close() error
}
