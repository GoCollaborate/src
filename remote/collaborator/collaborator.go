package collaborator

import (
	"fmt"
	"github.com/GoCollaborate/cmd"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote/remoteshared"
	"github.com/GoCollaborate/server"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/store"
	"github.com/GoCollaborate/utils"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

// Collaborator is a helper struct to operate on any inner trxs
type Collaborator struct {
	CardCase  Case
	Publisher *server.Publisher
	Workable  server.Workable
	Logger    *logger.Logger
}

func NewCollaborator(pbls *server.Publisher, localLogger *logger.Logger) *Collaborator {
	return &Collaborator{*newCase(), pbls, server.Dummy(), localLogger}
}

func newCase() *Case {
	return &Case{cmd.Vars().CaseID, Exposed{make(map[string]remoteshared.Card), time.Now().Unix()}, Reserved{*remoteshared.Default(), remoteshared.Card{}}}
}

func (clbt *Collaborator) Hire(wk server.Workable) {
	clbt.Workable = wk

	err := populate(&clbt.CardCase)
	if err != nil {
		panic(err)
	}

	logger.LogNormal("Case Identifier:")
	logger.LogListPoint(clbt.CardCase.CaseID)
	logger.LogNormal("Local Address:")
	ip := utils.GetLocalIP()
	logger.LogListPoint(ip)

	logger.LogNormal("Cards:")

	cp := clbt.CardCase.Cards()

	for _, h := range cp {
		var alive string
		if h.Alive {
			alive = "Alive"
		} else {
			alive = "Terminated"
		}
		logger.LogListPoint(h.GetFullIP(), alive)
	}

	local := clbt.CardCase.Local

	if card, ok := clbt.CardCase.Exposed.Cards[local.GetFullIP()]; !ok || !card.IsEqualTo(&local) {

		clbt.CardCase.Exposed.Cards[local.GetFullIP()] = local
		clbt.CardCase.Stamp()
		clbt.CardCase.writeStream()
	}

	clbt.launchServer()

	go func() {
		// start looking for other peer Collaborators
		clbt.Catchup()

		for {
			<-time.After(constants.DefaultSynInterval)
			err := populate(&clbt.CardCase)
			if err != nil {
				panic(err)
				return
			}
			clbt.Catchup()
			// clean collaborators that are no longer alive
			clbt.Clean()
		}
	}()

	/*
		stop the rpc connection at any time by calling Disconnect() function
	*/
	// Case.Disconnect()
}

// bridge to peer servers
func (clbt *Collaborator) Catchup() {
	c := &clbt.CardCase

	cp := c.Cards()

	max := cmd.Vars().GossipNum
	done := 0

	for key, e := range cp {
		// stop gossipping if already walk through max nodes
		if done > max {
			break
		}

		if e.Alive && (!e.IsEqualTo(&c.Local)) {
			client, err := launchClient(e.IP, e.Port)

			if err != nil {
				logger.LogWarning("Connection failed while bridging")
				override := c.Exposed.Cards[key]
				override.Alive = false
				c.Exposed.Cards[key] = override
				c.writeStream()
				continue
			}

			from := c.Local.GetFullExposureCard()
			to := e
			var in *remoteshared.CardMessage = remoteshared.NewCardMessageWithOptions(c.Cluster(), from, to, c.Cards(), c.TimeStamp())
			var out *remoteshared.CardMessage = remoteshared.NewCardMessage()
			err = client.Exchange(in, out)

			if err != nil {
				logger.LogWarning("Calling method failed while bridging")
				continue
			}

			if Compare(c, out) {
				cs, _ := Merge(c, out)
				cs.writeStream()
			}
		}

		done++
	}
}

func (clbt *Collaborator) Clean() {
	cp := clbt.CardCase.Cards()
	for k, c := range cp {
		if !c.Alive {
			clbt.CardCase.Terminate(k)
		}
	}
	clbt.CardCase.writeStream()

	cp = clbt.CardCase.Cards()

	logger.LogNormal("Cards:")
	for _, h := range cp {
		var alive string
		if h.Alive {
			alive = "Alive"
		} else {
			alive = "Terminated"
		}
		logger.LogListPoint(h.GetFullIP(), alive)
	}
}

func (clbt *Collaborator) Handle(router *mux.Router) *mux.Router {
	// register tasks existing in publisher
	// reflect api endpoint based on exposed task (function) name
	// once api get called, use distribute
	fs := store.GetInstance()
	for _, jobFunc := range fs.SharedJobs {
		logger.LogNormalWithPrefix("Job Linked: ", jobFunc.Signature)
		// shared registry
		clbt.HandleShared(router, jobFunc.Signature, jobFunc)
	}

	for _, jobFunc := range fs.LocalJobs {
		logger.LogNormalWithPrefix("Job Linked: ", jobFunc.Signature)
		// local registry
		clbt.HandleLocal(router, jobFunc.Signature, jobFunc)
	}
	return router
}

func (clbt *Collaborator) Distribute(sources ...*task.Task) ([]*task.Task, error) {
	b := clbt.CardCase
	var result []*task.Task
	l1 := int64(len(sources))
	l2 := int64(len(b.Cards()))
	if l2 < 1 {
		return []*task.Task{}, constants.ErrNoPeers
	}
	i := int64(-1)
	for _, e := range b.Cards() {
		i++
		var l = l1 % l2
		if (i * l) >= l1 {
			break
		}
		s := sources[(i * l):((i + 1) * l)]
		if e.IsEqualTo(&b.Local) {
			clbt.Publisher.LocalDistribute(s...)
			continue
		}
		client, err := launchClient(e.IP, e.Port)
		if err != nil {
			logger.LogWarning("Connection failed while connecting")
			continue
		}
		err = client.Distribute(s, &result)
		if err != nil {
			logger.LogWarning("Calling method failed while terminating")
			continue
		}
	}
	return result, nil
}

func (clbt *Collaborator) SyncDistribute(sources map[int64]*task.Task) (map[int64]*task.Task, error) {
	b := clbt.CardCase
	l1 := int64(len(sources))
	l2 := int64(len(b.Cards()))
	var (
		result  map[int64]*task.Task = make(map[int64]*task.Task)
		counter int64                = 0
	)

	// task channel
	chs := make(map[int64]chan *task.Task)

	if l2 < 1 {
		return result, constants.ErrNoPeers
	}
	// waiting for rpc responses
	printProgress(counter, l1)
	i := int64(-1)
	for _, k := range sources {
		i++
		var l = (l1 + i - 1) % l2
		if (i * l) >= l1 {
			break
		}
		e := b.ReturnByPosInt64(l)

		// publish to local if local Card/Card is down
		if e.IsEqualTo(&b.Local) || !e.Alive {
			// copy a pointer for concurrent map access
			p := *k
			chs[i] = clbt.Publisher.SyncDistribute(&p)
			continue
		}

		// publish to remote
		client, err := launchClient(e.IP, e.Port)
		if err != nil {
			// re-publish to local if failed
			logger.LogWarning("Connection failed while connecting")
			logger.LogWarning("Republish task to local")
			// copy a pointer for concurrent map access
			p := *k
			chs[i] = clbt.Publisher.SyncDistribute(&p)
			continue
		}
		// copy a pointer for concurrent map access
		s := make(map[int64]*task.Task)
		s[i] = k
		chs[i] = client.SyncDistribute(&s, &s)
	}
	// Wait for responses
	for {
		for i, ch := range chs {
			select {
			case t := <-ch:
				counter++
				result[int64(i)] = t
				printProgress(counter, l1)
			}
		}
		if counter >= l1 {
			break
		}
	}

	logger.LogProgress("All task responses collected")
	return result, nil
}

func registerRemote(server *rpc.Server, remote RemoteMethods) {
	server.RegisterName("RemoteMethods", remote)
}

func (c *Collaborator) launchServer() {
	go func() {
		methods := NewLocalMethods(c.Workable)
		server := rpc.NewServer()
		registerRemote(server, methods)
		// debug setup for regcenter to check services alive
		server.HandleHTTP("/", "debug")
		l, e := net.Listen("tcp", ":"+strconv.Itoa(c.CardCase.Local.Port))
		if e != nil {
			logger.LogErrorWithPrefix("Listen Error:", e)
		}
		http.Serve(l, nil)
	}()
	return
}

func launchClient(endpoint string, port int) (*RPCClient, error) {
	clientContact := remoteshared.Card{endpoint, port, true, ""}
	client, err := rpc.DialHTTP("tcp", clientContact.GetFullIP())
	if err != nil {
		logger.LogErrorWithPrefix("Dialing:", err)
		return &RPCClient{}, err
	}
	Card := &RPCClient{Client: client}
	return Card, nil
}

func (clbt *Collaborator) HandleLocal(router *mux.Router, api string, jobFunc *store.JobFunc) {
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		job := jobFunc.F(w, r)
		for s := job.Front(); s != nil; s = s.Next() {
			clbt.Publisher.LocalDistribute(s.TaskSet...)
		}
	}).Methods(jobFunc.Methods...)
}

func (clbt *Collaborator) HandleShared(router *mux.Router, api string, jobFunc *store.JobFunc) {
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		job := jobFunc.F(w, r)
		for s := job.Front(); s != nil; s = s.Next() {
			clbt.Publisher.SharedDistribute(s.TaskSet...)
		}
	}).Methods(jobFunc.Methods...)
}

func printProgress(num interface{}, tol interface{}) {
	logger.LogProgress("Waiting for task responses[" + fmt.Sprintf("%v", num) + "/" + fmt.Sprintf("%v", tol) + "]")
}
