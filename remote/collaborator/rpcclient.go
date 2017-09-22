package collaborator

import (
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote/remoteshared"
	"github.com/GoCollaborate/server/task"
	"net/rpc"
)

type RemoteMethods interface {
	Exchange(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error
	Disconnect(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error
	Distribute(source []*task.Task, result *[]*task.Task) error
	SyncDistribute(source *map[int64]*task.Task, result *map[int64]*task.Task) error
}

type RPCClient struct {
	Client *rpc.Client
}

func (c *RPCClient) Exchange(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error {
	err := c.Client.Call("RemoteMethods.Exchange", in, out)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}

func (c *RPCClient) Disconnect(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error {
	err := c.Client.Call("RemoteMethods.Disconnect", in, out)
	if err != nil {
		logger.LogError(err.Error())
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
