package remoteshared

import (
	"github.com/GoCollaborate/server/task"
)

type Collaboratable interface {
	SyncDistribute(sources []*task.Task) ([]*task.Task, error)
}
