package collaborator

import (
	"github.com/GoCollaborate/src/artifacts/card"
	"github.com/GoCollaborate/src/artifacts/message"
	"github.com/GoCollaborate/src/artifacts/task"
	. "github.com/GoCollaborate/src/collaborator/services"
	"github.com/GoCollaborate/src/constants"
	"github.com/GoCollaborate/src/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ServiceClientStub struct {
	RPCServiceClient
}

func NewServiceClientStub(endpoint string, port int32, secure ...bool) (stub *ServiceClientStub, err error) {
	if len(secure) < 1 {
		clientContact := card.Card{endpoint, port, true, "", false}

		conn, err := grpc.Dial(
			clientContact.GetFullIP(),
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(constants.DEFAULT_RPC_DIAL_TIMEOUT),
		)

		if err != nil {
			logger.LogError("Dialing:", err)
			logger.GetLoggerInstance().LogError("Dialing:", err)
			return &ServiceClientStub{}, err
		}

		return &ServiceClientStub{NewRPCServiceClient(conn)}, nil
	}
	// todo: change return to TLS client
	return &ServiceClientStub{}, nil
}

func (stub *ServiceClientStub) DistributeAsync(source *map[int]*task.Task) chan *task.Task {
	ch := make(chan *task.Task)

	go func() {
		defer close(ch)
		enc, _ := Encode(source)
		dec, err := stub.RPCServiceClient.Distribute(context.Background(), enc)
		if err != nil {
			logger.LogError("Connection Error:" + err.Error())
			return
		}
		result, _ := Decode(dec)
		for _, t := range *result {
			ch <- t
		}
	}()

	return ch
}

func (stub *ServiceClientStub) Exchange(in *message.CardMessage) (*message.CardMessage, error) {
	return stub.RPCServiceClient.Exchange(context.Background(), in)
}
