package remote

import (
	"fmt"
	"net/rpc"
)

type RemoteMethods interface {
	Signal(arg *ContactBook, update *ContactBook) error
	Disconnect(arg *ContactBook, update *ContactBook) error
	Terminate(arg *ContactBook, update *ContactBook) error
}

type RPCClient struct {
	Client *rpc.Client
}

func (c *RPCClient) Signal(arg *ContactBook, update *ContactBook) error {
	err := c.Client.Call("RemoteMethods.Signal", arg, update)
	if err != nil {
		fmt.Println("Connection Error:", err)
		return err
	}
	return nil
}

func (c *RPCClient) Disconnect(arg *ContactBook, update *ContactBook) error {
	err := c.Client.Call("RemoteMethods.Disconnect", arg, update)
	fmt.Println(update)
	if err != nil {
		fmt.Println("Connection Error:", err)
		return err
	}
	return nil
}

func (c *RPCClient) Terminate(arg *ContactBook, update *ContactBook) error {
	err := c.Client.Call("RemoteMethods.Terminate", arg, update)
	if err != nil {
		fmt.Println("Connection Error:", err)
		return err
	}
	return nil
}
