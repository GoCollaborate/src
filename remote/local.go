package remote

import (
	"github.com/GoCollaborate/logger"
)

type LocalMethods int

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
