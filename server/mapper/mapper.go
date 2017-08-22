package mapper

import (
	"github.com/GoCollaborate/server/task"
)

type Mapper interface {
	Map(t task.Task) (map[int64]task.Task, error)
}
