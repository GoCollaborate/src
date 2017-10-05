package collaborator

import (
	"github.com/GoCollaborate/artifacts/iworkable"
	"github.com/GoCollaborate/artifacts/message"
	"github.com/GoCollaborate/artifacts/task"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/wrappers/messageHelper"
)

type LocalMethods struct {
	workable iworkable.Workable
}

func NewLocalMethods(wk iworkable.Workable) *LocalMethods {
	return &LocalMethods{wk}
}

func (l *LocalMethods) Exchange(in *message.CardMessage, out *message.CardMessage) error {
	logger.LogNormal("Card message from another Collaborator received")
	err := messageHelper.Exchange(in, out)
	return err
}

func (l *LocalMethods) DistributeSync(source *map[int]*task.Task, result *map[int]*task.Task) error {
	logger.LogNormal("Task from another Collaborator received")
	s := *source
	err := l.workable.DoneMulti(s)
	*result = s
	return err
}
