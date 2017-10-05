package iremote

import (
	"github.com/GoCollaborate/artifacts/task"
)

type ICollaboratable interface {
	SyncDistribute(sources []*task.Task) ([]*task.Task, error)
}
