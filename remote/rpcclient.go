package remote

import (
	"fmt"
	"net/rpc"
)

type RemoteMethods interface {
	Signal(arg *Config, update *Config) error
	Disconnect(arg *Config, update *Config) error
	Terminate(arg *Config, update *Config) error
}

type RPCClient struct {
	Client *rpc.Client
}

func (c *RPCClient) Signal(arg *Config, update *Config) error {
	err := c.Client.Call("RemoteMethods.Signal", arg, update)
	if err != nil {
		fmt.Println("Connection Error:", err)
		return err
	}
	return nil
}

func (c *RPCClient) Disconnect(arg *Config, update *Config) error {
	err := c.Client.Call("RemoteMethods.Disconnect", arg, update)
	fmt.Println(update)
	if err != nil {
		fmt.Println("Connection Error:", err)
		return err
	}
	return nil
}

func (c *RPCClient) Terminate(arg *Config, update *Config) error {
	err := c.Client.Call("RemoteMethods.Terminate", arg, update)
	if err != nil {
		fmt.Println("Connection Error:", err)
		return err
	}
	return nil
}
