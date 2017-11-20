package collaborator

import (
	"github.com/GoCollaborate/src/artifacts/iworkable"
	"github.com/GoCollaborate/src/artifacts/message"
	"github.com/GoCollaborate/src/artifacts/task"
	. "github.com/GoCollaborate/src/collaborator/services"
	"github.com/GoCollaborate/src/logger"
	"github.com/GoCollaborate/src/wrappers/messageHelper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

func NewServiceServerStub(wk iworkable.Workable) *ServiceServerStub {
	return &ServiceServerStub{wk}
}

type ServiceServerStub struct {
	workable iworkable.Workable
}

func LaunchServer(addr string, wkb iworkable.Workable) {
	go func() {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		RegisterRPCServiceServer(grpcServer, &ServiceServerStub{wkb})
		grpcServer.Serve(lis)
	}()
}

func (stub *ServiceServerStub) Exchange(
	ctx context.Context,
	in *message.CardMessage) (*message.CardMessage, error) {

	var (
		out = message.NewCardMessage()
		err error
	)

	logger.LogNormal("Card message from another Collaborator received")
	out, err = messageHelper.Exchange(in)
	return out, err
}

// Local implementation of rpc method Distribute()
func (stub *ServiceServerStub) Distribute(
	ctx context.Context,
	in *task.TaskPayload) (out *task.TaskPayload, err error) {

	var (
		src *map[int]*task.Task
		rst *map[int]*task.Task
	)

	src, err = Decode(in)
	rst, err = stub.distribute(src)
	out, err = Encode(rst)

	return out, err
}

// Local implementation of rpc method Distribute()
func (stub *ServiceServerStub) distribute(
	source *map[int]*task.Task) (*map[int]*task.Task, error) {
	var (
		result = new(map[int]*task.Task)
		err    error
	)
	logger.LogNormal("Task from another Collaborator received")
	s := *source
	err = stub.workable.DoneMulti(s)
	*result = s
	return result, err
}
