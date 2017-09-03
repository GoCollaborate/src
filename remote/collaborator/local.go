package collaborator

import (
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server"
	"github.com/GoCollaborate/server/task"
)

type LocalMethods struct {
	workable server.Workable
}

func NewLocalMethods(wk server.Workable) *LocalMethods {
	return &LocalMethods{wk}
}

func (l *LocalMethods) Signal(arg *ContactBook, update *ContactBook) error {
	logger.LogNormal("Signal Acquired!")
	upd, err := arg.RemoteLoad()

	// update local config to remote call
	update.Agents = upd.Agents
	update.Local = upd.Local
	update.TimeStamp = upd.TimeStamp
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalMethods) Disconnect(arg *ContactBook, update *ContactBook) error {
	logger.LogNormal("Disconnect:" + arg.Local.GetFullIP())
	upd, err := arg.RemoteDisconnect()

	// update local config to remote call
	update.Agents = upd.Agents
	update.Local = upd.Local
	update.TimeStamp = upd.TimeStamp

	if err != nil {
		return err
	}
	return nil
}

func (l *LocalMethods) Terminate(arg *ContactBook, update *ContactBook) error {
	logger.LogNormal("Terminate:" + arg.Local.GetFullIP())
	upd, err := arg.RemoteTerminate()

	// update local config to remote call
	update.Agents = upd.Agents
	upd.Local.Alive = false
	update.Local = upd.Local
	update.TimeStamp = upd.TimeStamp

	if err != nil {
		return err
	}
	return nil
}

func (l *LocalMethods) Distribute(source []*task.Task, result *[]*task.Task) error {
	logger.LogNormal("Distribute...")
	l.workable.Enqueue(source...)
	result = &source
	return nil
}

func (l *LocalMethods) SyncDistribute(source *map[int64]*task.Task, result *map[int64]*task.Task) error {
	logger.LogNormal("SyncDistribute...")
	s := *source
	r := make(map[int64]*task.Task)

	for i, m := range s {
		err := l.workable.Done(m)
		if err != nil {
			return err
		}
		r[i] = m
	}
	*result = r
	return nil
}
