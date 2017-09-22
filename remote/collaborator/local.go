package collaborator

import (
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote/remoteshared"
	"github.com/GoCollaborate/server"
	"github.com/GoCollaborate/server/task"
)

type LocalMethods struct {
	workable server.Workable
}

func NewLocalMethods(wk server.Workable) *LocalMethods {
	return &LocalMethods{wk}
}

func (l *LocalMethods) Exchange(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error {
	logger.LogNormal("Card message from another Collaborator received")

	err := RemoteLoad(in, out)

	if err != nil {
		return err
	}
	return nil
}

func (l *LocalMethods) Disconnect(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error {
	from := in.From()
	logger.LogNormal("Disconnect:" + from.GetFullIP())

	err := RemoteDisconnect(in, out)

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
