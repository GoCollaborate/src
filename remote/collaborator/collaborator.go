package collaborator

import (
	"encoding/json"
	"fmt"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote/remoteshared"
	"github.com/GoCollaborate/server"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

type BookKeeper struct {
	Book      ContactBook
	Publisher *server.Publisher `json:"-"`
	Workable  server.Workable   `json:"-"`
	Logger    *logger.Logger    `json:"-"`
}

func NewBookKeeper(pbls *server.Publisher, localLogger *logger.Logger) *BookKeeper {
	return &BookKeeper{ContactBook{}, pbls, server.Dummy(), localLogger}
}

func (bk *BookKeeper) NewBook() {
	bk.Book = ContactBook{[]remoteshared.Agent{}, remoteshared.Agent{}, *Default(), false, false, time.Now().Unix()}
}

func (bk *BookKeeper) Watch(wk server.Workable) {
	bk.Workable = wk
	bk.Book.Load()
	bk.LaunchServer()
	// bridge to remote servers
	bk.Book.Bridge()

	/*
		terminate the rpc connection at any time by calling Terminate() function
	*/
	// contactBook.Terminate()
}

func (bk *BookKeeper) Handle(router *mux.Router) *mux.Router {
	// register tasks existing in publisher
	// reflect api endpoint based on exposed task (function) name
	// once api get called, use distribute
	for _, tskFunc := range bk.Publisher.SharedTasks {
		// tskFunc is an iterator, and it shouldn't
		// be accessed during the runtime
		api := utils.StripRouteToAPIRoute(utils.ReflectFuncName(tskFunc.F))
		logger.LogNormalWithPrefix(logger.NORMAL, "Task Linked: ", api)
		// shared registry
		bk.Publisher.HandleShared(router, api, tskFunc)
	}

	for _, tskFunc := range bk.Publisher.LocalTasks {
		// tskFunc is an iterator, and it shouldn't
		// be accessed during the runtime
		api := utils.StripRouteToAPIRoute(utils.ReflectFuncName(tskFunc.F))
		logger.LogNormalWithPrefix(logger.NORMAL, "Task Linked: ", api)
		// local registry
		bk.Publisher.HandleLocal(router, api, tskFunc)
	}
	router.HandleFunc("/", handlerFuncBookKeeperEntry)
	return router
}

// current RPC port
func Default() *remoteshared.Agent {
	return &remoteshared.Agent{"localhost", GetPort(), true, ""}
}

// local is the local config of server
type ContactBook struct {
	Agents                 []remoteshared.Agent `json:"agents,omitempty"`
	Local                  remoteshared.Agent   `json:"local,omitempty"`
	Coordinator            remoteshared.Agent   `json:"coordinator,omitempty"`
	ForbidPointToPointConn bool                 `json:"forbidPointToPointConn,omitempty"`
	IsCoordinator          bool                 `json:"isCoordinator,omitempty"`
	TimeStamp              int64                `json:"timestamp,omitempty"`
}

func Populate(cfg *ContactBook) error {
	bytes, err := ioutil.ReadFile(constants.DefaultContactBookPath)
	if err != nil {
		panic(err)
	}
	// unmarshal, overwrite default if already existed in config file
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}

func (c *ContactBook) Load() (*ContactBook, error) {
	err := Populate(c)
	if err != nil {
		return c, err
	}

	logger.LogNormal("Localaddress:")
	ip := GetLocalIP()
	logger.LogNormal(ip)

	var exist bool = false
	var idx int

	logger.LogNormal("Successfully load config file...")
	logger.LogNormal("Hosts:")

	for i, h := range c.Agents {
		if h.IsEqualTo(&c.Local) {
			exist = true
			idx = i
		}
		logger.LogNormal(h.GetFullIP())
	}

	// update if not exist
	if !exist {
		c.Agents = append(c.Agents, c.Local)
		c.Sync()
	} else {
		// activate if not alive
		if !c.Agents[idx].Alive {
			c.Agents[idx] = remoteshared.Agent{c.Agents[idx].IP, c.Agents[idx].Port, true, ""}
			c.Sync()
		}
	}
	return c, nil
}

func (c *ContactBook) RemoteLoad() (*ContactBook, error) {
	var localBook *ContactBook = new(ContactBook)
	err := Populate(localBook)
	if err != nil {
		return c, err
	}
	var exist bool = false
	var update bool = false

	update = Compare(localBook, c)

	for _, h := range c.Agents {
		if h.IsEqualTo(&localBook.Local) {
			exist = true
		}
	}
	// update if not exist
	if !exist {
		c.Agents = append(c.Agents, c.Local)
		localBook.Agents = append(c.Agents, c.Local)
		update = true
	}
	if update {
		localBook.Sync()
	}
	return c, nil
}

func (c *ContactBook) RemoteDisconnect() (*ContactBook, error) {
	var localBook *ContactBook = new(ContactBook)
	err := Populate(localBook)
	if err != nil {
		return c, err
	}
	var index = -1
	var update bool = false

	update = Compare(localBook, c)

	for i, h := range c.Agents {
		if h.IsEqualTo(&c.Local) {
			index = i
		}
	}
	// update if not exist
	if index < 0 {
		c.Agents = append(c.Agents, c.Local)
		localBook.Agents = append(c.Agents, c.Local)
		update = true
	} else {
		// deactivate if alive
		if c.Agents[index].Alive {
			c.Agents[index] = remoteshared.Agent{c.Agents[index].IP, c.Agents[index].Port, false, ""}
			c.Stamp()
			for _, a := range localBook.Agents {
				if a.IsEqualTo(&c.Agents[index]) {
					a.Alive = false
				}
			}
			update = true
		}
	}
	if update {
		localBook.Sync()
	}
	return c, nil
}

func (c *ContactBook) RemoteTerminate() (*ContactBook, error) {
	var localBook *ContactBook = new(ContactBook)
	err := Populate(localBook)
	if err != nil {
		return c, err
	}
	var index int = -1
	var update bool = false

	for i, h := range localBook.Agents {
		if h.IsEqualTo(&c.Local) {
			index = i
			update = true
		}
	}
	// update if exist
	if index >= 0 {
		localBook.RemoveAgent(index)
		index = -1
	}
	for j, h := range c.Agents {
		if h.IsEqualTo(&c.Local) {
			index = j
		}
	}
	if index >= 0 {
		c.RemoveAgent(index)
		c.Stamp()
	}

	if update {
		localBook.Sync()
	}
	return c, nil
}

// bridge to peer servers
func (c *ContactBook) Bridge() {
	for _, e := range c.Agents {
		if e.Alive && (!e.IsEqualTo(&c.Local)) {
			client, err := LaunchClient(e.IP, e.Port)
			if err != nil {
				logger.LogWarning("Connection failed while bridging...")
				continue
			}
			var u ContactBook
			err = client.Signal(c, &u)
			if err != nil {
				logger.LogWarning("Calling method failed while bridging...")
				continue
			}
			u.Update(c).WriteStream()
		}
	}
}

func (c *ContactBook) Disconnect() {
	for _, e := range c.Agents {
		if e.Alive && (!e.IsEqualTo(&c.Local)) {
			client, err := LaunchClient(e.IP, e.Port)
			if err != nil {
				logger.LogWarning("Connection failed while disconnecting...")
				continue
			}
			var u ContactBook
			err = client.Disconnect(c, &u)
			if err != nil {
				logger.LogWarning("Calling method failed while disconnecting...")
				continue
			}
			u.Update(c).WriteStream()
		}
	}
}

func (c *ContactBook) Terminate() {
	for _, e := range c.Agents {
		if e.Alive && (!e.IsEqualTo(&c.Local)) {
			client, err := LaunchClient(e.IP, e.Port)
			if err != nil {
				logger.LogWarning("Connection failed while terminating...")
				continue
			}
			var u ContactBook
			err = client.Terminate(c, &u)
			if err != nil {
				logger.LogWarning("Calling method failed while terminating...")
				continue
			}
			u.Update(c).WriteStream()
		}
	}
}

// Sync will also update TimeStamp
func (c *ContactBook) Sync() {
	c.TimeStamp = time.Now().Unix()
	mal, err := json.Marshal(&c)
	err = ioutil.WriteFile(constants.DefaultContactBookPath, mal, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
}

// Update() does not update TimeStamp
func (c *ContactBook) Update(update *ContactBook) (cfg *ContactBook) {
	Compare(update, c)
	return c
}

func (c *ContactBook) WriteStream() {
	mal, err := json.Marshal(&c)
	err = ioutil.WriteFile(constants.DefaultContactBookPath, mal, os.ModeExclusive)
	if err != nil {
		logger.LogError(err)
	}
}

func (c *ContactBook) RemoveAgent(index int) *ContactBook {
	c.Agents = append(c.Agents[:(index)], c.Agents[(index+1):]...)
	return c
}

func (c *ContactBook) Stamp() *ContactBook {
	c.TimeStamp = time.Now().Unix()
	return c
}

func (bk *BookKeeper) Distribute(sources ...task.Task) ([]task.Task, error) {
	c := bk.Book
	var result []task.Task
	l1 := len(sources)
	l2 := len(c.Agents)
	if l2 < 1 {
		return []task.Task{}, constants.ErrNoPeers
	}
	for i, e := range c.Agents {
		var l = l1 % l2
		if (i * l) >= l1 {
			break
		}
		s := sources[(i * l):((i + 1) * l)]
		if e.GetFullIP() == c.Local.GetFullIP() {
			bk.Publisher.LocalDistribute(s...)
			continue
		}
		client, err := LaunchClient(e.IP, e.Port)
		if err != nil {
			logger.LogWarning("Connection failed while connecting...")
			continue
		}
		err = client.Distribute(&s, &result)
		if err != nil {
			logger.LogWarning("Calling method failed while terminating...")
			continue
		}
	}
	return result, nil
}

func (bk *BookKeeper) SyncDistribute(sources map[int64]task.Task) (map[int64]task.Task, error) {
	var (
		result  map[int64]task.Task
		counter int64 = 0
	)

	c := bk.Book
	l1 := int64(len(sources))
	l2 := int64(len(c.Agents))

	// task channel
	chs := make([]chan error, l1)

	if l2 < 1 {
		return result, constants.ErrNoPeers
	}
	// waiting for rpc responses
	printProgress(counter, l1)
	for i, k := range sources {
		var l = (l1 + i - 1) % l2
		if (i * l) >= l1 {
			break
		}
		e := c.Agents[l]

		// publish to local if local agent/agent is down
		if e.IsEqualTo(&c.Local) || !e.Alive {
			chs[i] = bk.Publisher.SyncDistribute(k)
			continue
		}

		// publish to remote
		client, err := LaunchClient(e.IP, e.Port)
		if err != nil {
			// re-publish to local if failed
			logger.LogWarning("Connection failed while connecting...")
			c.Agents[l].Alive = false
			logger.LogWarning("Republish task to local...")
			chs[i] = bk.Publisher.SyncDistribute(k)
			continue
		}
		s := make(map[int64]task.Task)
		s[i] = k
		chs[i] = client.SyncDistribute(&s, &result)
	}
	for {
		for _, ch := range chs {
			select {
			case <-ch:
				counter++
				printProgress(counter, l1)
			}
		}
		if counter >= l1 {
			break
		}
	}

	logger.LogProgress("All task responses collected...")
	return result, nil
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback then display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func Compare(a *ContactBook, b *ContactBook) bool {
	if a.TimeStamp < b.TimeStamp {
		a.Agents = b.Agents
		a.TimeStamp = b.TimeStamp
		return true
	}
	b.Agents = a.Agents
	b.TimeStamp = a.TimeStamp
	return false
}

func RegisterRemote(server *rpc.Server, remote RemoteMethods) {
	server.RegisterName("RemoteMethods", remote)
}

func (c *BookKeeper) LaunchServer() {
	go func() {
		methods := NewLocalMethods(c.Workable)
		server := rpc.NewServer()
		RegisterRemote(server, methods)
		// debug setup for regcenter to check services alive
		server.HandleHTTP("/", "debug")
		l, e := net.Listen("tcp", ":"+strconv.Itoa(c.Book.Local.Port))
		if e != nil {
			logger.LogErrorWithPrefix(logger.NORMAL, "Listen Error:", e)
		}
		http.Serve(l, nil)
	}()
	return
}

func LaunchClient(endpoint string, port int) (*RPCClient, error) {
	clientContact := remoteshared.Agent{endpoint, port, true, ""}
	client, err := rpc.DialHTTP("tcp", clientContact.GetFullIP())
	if err != nil {
		logger.LogErrorWithPrefix(logger.NORMAL, "Dialing:", err)
		return &RPCClient{}, err
	}
	agent := &RPCClient{Client: client}
	return agent, nil
}

func handlerFuncBookKeeperEntry(w http.ResponseWriter, r *http.Request) {
	logger.LogNormal("Bookkeeper Println...")
}

func printProgress(num interface{}, tol interface{}) {
	logger.LogProgress("Waiting for task responses[" + fmt.Sprintf("%v", num) + "/" + fmt.Sprintf("%v", tol) + "]...")
}
