package collaborator

import (
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/utils"
	"net/rpc"
)

type RemoteMethods interface {
	Signal(arg *ContactBook, update *ContactBook) error
	Disconnect(arg *ContactBook, update *ContactBook) error
	Terminate(arg *ContactBook, update *ContactBook) error
	Distribute(source []*task.Task, result *[]*task.Task) error
	SyncDistribute(source *map[int64]*task.Task, result *map[int64]*task.Task) error
}

type RPCClient struct {
	Client *rpc.Client
}

func (c *RPCClient) Signal(arg *ContactBook, update *ContactBook) error {
	arg.Local.IP = utils.MapToExposureAddress(arg.Local.IP)
	err := c.Client.Call("RemoteMethods.Signal", arg, update)
	update.Local.IP = utils.MapToExposureAddress(update.Local.IP)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}

func (c *RPCClient) Disconnect(arg *ContactBook, update *ContactBook) error {
	arg.Local.IP = utils.MapToExposureAddress(arg.Local.IP)
	err := c.Client.Call("RemoteMethods.Disconnect", arg, update)
	update.Local.IP = utils.MapToExposureAddress(update.Local.IP)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}

func (c *RPCClient) Terminate(arg *ContactBook, update *ContactBook) error {
	arg.Local.IP = utils.MapToExposureAddress(arg.Local.IP)
	err := c.Client.Call("RemoteMethods.Terminate", arg, update)
	update.Local.IP = utils.MapToExposureAddress(update.Local.IP)
	if err != nil {
		logger.LogError("Connection Error:" + err.Error())
		return err
	}
	return nil
}

func (c *RPCClient) Distribute(source []*task.Task, result *[]*task.Task) error {
	go func() {
		err := c.Client.Call("RemoteMethods.Distribute", source, result)
		if err != nil {
			logger.LogError("Connection Error:" + err.Error())
		}
	}()
	return nil
}

func (c *RPCClient) SyncDistribute(source *map[int64]*task.Task, result *map[int64]*task.Task) chan *task.Task {
	ch := make(chan *task.Task)
	go func() {
		defer close(ch)
		err := c.Client.Call("RemoteMethods.SyncDistribute", source, result)
		if err != nil {
			logger.LogError("Connection Error:" + err.Error())
		}
		for _, t := range *result {
			ch <- t
		}
	}()
	return ch
}
