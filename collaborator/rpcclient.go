package collaborator

import (
	"github.com/GoCollaborate/src/artifacts/message"
	"github.com/GoCollaborate/src/artifacts/task"
	"github.com/GoCollaborate/src/logger"
	"net/rpc"
)

type RemoteMethods interface {
	Exchange(in *message.CardMessage, out *message.CardMessage) error
	DistributeSync(source *map[int]*task.Task, result *map[int]*task.Task) error
}

type RPCClient struct {
	Client *rpc.Client
}

func (c *RPCClient) Exchange(in *message.CardMessage, out *message.CardMessage) error {
	err := c.Client.Call("RemoteMethods.Exchange", in, out)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}

func (c *RPCClient) DistributeSync(source *map[int]*task.Task, result *map[int]*task.Task) chan *task.Task {
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
