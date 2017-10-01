package collaborator

import (
	"fmt"
	"github.com/GoCollaborate/cmd"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote/remoteshared"
	"github.com/GoCollaborate/server"
	"github.com/GoCollaborate/server/executor"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/store"
	"github.com/GoCollaborate/utils"
	"github.com/GoCollaborate/web"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

const (
	MAXIDLECONNECTIONS int           = 20
	REQUESTTIMEOUT     time.Duration = 5
	UPDATEINTERVAL     time.Duration = 1
)

// Collaborator is a helper struct to operate on any inner trxs
type Collaborator struct {
	CardCase Case
	Workable server.Workable
}

func NewCollaborator() *Collaborator {
	return &Collaborator{*newCase(), server.Dummy()}
}

func newCase() *Case {
	return &Case{cmd.Vars().CaseID, Exposed{make(map[string]remoteshared.Card), time.Now().Unix()}, Reserved{*remoteshared.Default(), remoteshared.Card{}}}
}

func (clbt *Collaborator) Join(wk server.Workable) {
	clbt.Workable = wk

	err := populate(&clbt.CardCase)
	if err != nil {
		panic(err)
	}

	logger.LogNormal("Case Identifier:")
	logger.GetLoggerInstance().LogNormal("Case Identifier:")
	logger.LogListPoint(clbt.CardCase.CaseID)
	logger.GetLoggerInstance().LogListPoint(clbt.CardCase.CaseID)

	ip := utils.GetLocalIP()
	logger.LogNormal("Local Address:")
	logger.LogListPoint(ip)
	logger.GetLoggerInstance().LogNormal("Local Address:")
	logger.GetLoggerInstance().LogListPoint(ip)

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
		logger.GetLoggerInstance().LogListPoint(h.GetFullIP(), alive)
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
				logger.GetLoggerInstance().LogWarning("Connection failed while bridging")
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
				logger.GetLoggerInstance().LogWarning("Calling method failed while bridging")
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

	// register dashboard
	router.HandleFunc("/", web.Index).Methods("GET").Name("Index")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(constants.LibUnixDir+"static/")))).Name("Static")
	router.HandleFunc("/dashboard/profile", web.Profile).Methods("GET").Name("Profile")
	router.HandleFunc("/dashboard/routes", web.Routes).Methods("GET").Name("Routes")
	router.HandleFunc("/dashboard/logs", web.Logs).Methods("GET").Name("Logs")

	// register tasks existing in publisher
	// reflect api endpoint based on exposed task (function) name
	// once api get called, use distribute
	fs := store.GetInstance()
	for _, jobFunc := range fs.SharedJobs {
		logger.LogNormal("Job Linked:")
		logger.LogListPoint(jobFunc.Signature)
		logger.GetLoggerInstance().LogNormal("Job Linked:")
		logger.GetLoggerInstance().LogListPoint(jobFunc.Signature)

		// shared registry
		clbt.HandleShared(router, jobFunc.Signature, jobFunc)
	}

	for _, jobFunc := range fs.LocalJobs {
		logger.LogNormal("Job Linked:")
		logger.LogListPoint(jobFunc.Signature)
		logger.GetLoggerInstance().LogNormal("Job Linked:")
		logger.GetLoggerInstance().LogListPoint(jobFunc.Signature)
		// local registry
		clbt.HandleLocal(router, jobFunc.Signature, jobFunc)
	}
	return router
}

func (clbt *Collaborator) SyncDistribute(sources map[int]*task.Task) (map[int]*task.Task, error) {
	cc := clbt.CardCase
	l1 := len(sources)
	l2 := len(cc.Cards())

	var (
		result  map[int]*task.Task = make(map[int]*task.Task)
		counter int                = 0
	)

	// task channel
	chs := make(map[int]chan *task.Task)

	if l2 < 1 {
		return result, constants.ErrNoPeers
	}
	// waiting for rpc responses
	printProgress(counter, l1)

	for k, v := range sources {
		var l = k % l2
		e := cc.ReturnByPos(l)

		// publish to local if local Card/Card is down
		if e.IsEqualTo(&cc.Local) || !e.Alive {
			// copy a pointer for concurrent map access
			p := *v
			chs[k] = clbt.DelayExecute(&p)
			continue
		}

		// publish to remote
		client, err := launchClient(e.IP, e.Port)
		if err != nil {
			// re-publish to local if failed
			logger.LogWarning("Connection failed while connecting")
			logger.LogWarning("Republish task to local")
			logger.GetLoggerInstance().LogWarning("Connection failed while connecting")
			logger.GetLoggerInstance().LogWarning("Republish task to local")
			// copy a pointer for concurrent map access
			p := *v
			chs[k] = clbt.DelayExecute(&p)
			continue
		}

		// copy a pointer for concurrent map access
		s := make(map[int]*task.Task)
		s[k] = v
		chs[k] = client.SyncDistribute(&s, &s)
	}

	// Wait for responses
	for {
		for i, ch := range chs {
			select {
			case t := <-ch:
				counter++
				result[i] = t
				printProgress(counter, l1)
			}
		}
		if counter >= l1 {
			break
		}
	}

	logger.LogProgress("All task responses collected")
	logger.GetLoggerInstance().LogProgress("All task responses collected")
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
			logger.LogError("Listen Error:", e)
			logger.GetLoggerInstance().LogError("Listen Error:", e)
		}
		http.Serve(l, nil)
	}()
	return
}

func launchClient(endpoint string, port int) (*RPCClient, error) {
	clientContact := remoteshared.Card{endpoint, port, true, ""}
	client, err := rpc.DialHTTP("tcp", clientContact.GetFullIP())
	if err != nil {
		logger.LogError("Dialing:", err)
		logger.GetLoggerInstance().LogError("Dialing:", err)
		return &RPCClient{}, err
	}
	Card := &RPCClient{Client: client}
	return Card, nil
}

func (clbt *Collaborator) HandleLocal(router *mux.Router, api string, jobFunc *store.JobFunc) {
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		job := jobFunc.F(w, r)
		var (
			counter = 0
		)
		for s := job.Front(); s != nil; s = s.Next() {
			exes, err := job.Exes(counter)
			if err != nil {
				logger.GetLoggerInstance().LogError(err)
				break
			}
			err = clbt.LocalDistribute(&s.TaskSet, exes)
			if err != nil {
				logger.GetLoggerInstance().LogError(err)
				break
			}
			counter++
		}
	}).Methods(jobFunc.Methods...).Name("Local Tasks")
}

func (clbt *Collaborator) HandleShared(router *mux.Router, api string, jobFunc *store.JobFunc) {
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		job := jobFunc.F(w, r)
		var (
			counter = 0
		)
		for s := job.Front(); s != nil; s = s.Next() {
			exes, err := job.Exes(counter)
			if err != nil {
				logger.GetLoggerInstance().LogError(err)
				break
			}

			err = clbt.SharedDistribute(&s.TaskSet, exes)
			if err != nil {
				logger.GetLoggerInstance().LogError(err)
				break
			}
			counter++
		}
	}).Methods(jobFunc.Methods...).Name("Shared Tasks")
}

func (clbt *Collaborator) LocalDistribute(pmaps *map[int]*task.Task, exe []string) error {
	clbt.Workable.Enqueue(*pmaps)
	return nil
}

// SharedDistribute will process the tasks of the same stage in memory
func (clbt *Collaborator) SharedDistribute(pmaps *map[int]*task.Task, stacks []string) error {
	var (
		err  error
		fs   = store.GetInstance()
		maps = *pmaps
	)
	for _, stack := range stacks {
		var (
			exe executor.Executor
		)
		exe, err = fs.GetExecutor(stack)

		if err != nil {
			return err
		}

		switch exe.Type() {
		// mapper
		case constants.ExecutorTypeMapper:
			maps, err = exe.Execute(maps)
			if err != nil {
				return err
			}
			maps, err = clbt.SyncDistribute(maps)
			// reducer
		case constants.ExecutorTypeReducer:
			maps, err = exe.Execute(maps)
			if err != nil {
				return err
			}
		default:
			maps, err = exe.Execute(maps)
			if err != nil {
				return err
			}
		}
	}
	*pmaps = maps

	return nil
}

func (clbt *Collaborator) DelayExecute(t *task.Task) chan *task.Task {
	ch := make(chan *task.Task)

	go func() {
		defer close(ch)
		err := clbt.Workable.Done(t)
		if err != nil {
			logger.LogError("Execution Error:" + err.Error())
			ch <- &task.Task{}
		}
		ch <- t
	}()
	return ch
}

func Delay(sec time.Duration) {
	tm := time.NewTimer(sec * time.Second)
	<-tm.C
}

func printProgress(num interface{}, tol interface{}) {
	logger.LogProgress("Waiting for task responses[" + fmt.Sprintf("%v", num) + "/" + fmt.Sprintf("%v", tol) + "]")
	logger.GetLoggerInstance().LogProgress("Waiting for task responses[" + fmt.Sprintf("%v", num) + "/" + fmt.Sprintf("%v", tol) + "]")
}
