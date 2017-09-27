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

func (l *LocalMethods) SyncDistribute(source *map[int]*task.Task, result *map[int]*task.Task) error {
	logger.LogNormal("Task from another Collaborator received")
	s := *source

	err := l.workable.DoneMulti(s)
	if err != nil {
		return err
	}

	*result = s
	return nil
}
