package imapper

import (
	"github.com/GoCollaborate/src/artifacts/task"
)

type IMapper interface {
	Map(sources map[int]*task.Task) (map[int]*task.Task, error)
}

func Default() *DefaultMapper {
	return new(DefaultMapper)
}

type DefaultMapper struct {
}

func (mp *DefaultMapper) Map(sources map[int]*task.Task) (map[int]*task.Task, error) {
	return map[int]*task.Task{}, nil
}
