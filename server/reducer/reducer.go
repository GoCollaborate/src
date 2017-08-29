package reducer

import (
	"github.com/GoCollaborate/server/task"
)

type Reducer interface {
	Reduce(sources map[int64]task.Task, result *task.Task) error
}

type DefaultReducer struct {
}

func Default() *DefaultReducer {
	return new(DefaultReducer)
}

func (rd *DefaultReducer) Reduce(sources map[int64]task.Task, result *task.Task) error {
	return nil
}
