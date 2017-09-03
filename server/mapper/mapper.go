package mapper

import (
	"github.com/GoCollaborate/server/task"
)

type Mapper interface {
	Map(t *task.Task) (map[int64]*task.Task, error)
}

func Default() *DefaultMapper {
	return new(DefaultMapper)
}

type DefaultMapper struct {
}

func (mp *DefaultMapper) Map(t *task.Task) (map[int64]*task.Task, error) {
	return map[int64]*task.Task{}, nil
}
