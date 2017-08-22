package remote

import (
	"encoding/json"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server"
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
	Book ContactBook
}

func NewBookKeeper() *BookKeeper {
	return new(BookKeeper)
}

func (bk *BookKeeper) WatchNewBook(pbls *server.Publisher, localLogger *logger.Logger) {
	bk.Book = ContactBook{[]Agent{}, Agent{}, *Default(), false, false, time.Now().Unix(), pbls, localLogger}
	bk.Book.Load()
	bk.Book.LaunchServer()
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
	c := &bk.Book
	for _, tskFunc := range c.Publisher.ExposedTasks {
		// tskFunc is an iterator, and it shouldn't
		// be accessed during the runtime
		api := utils.StripRouteToAPIRoute(utils.ReflectFuncName(tskFunc))
		logger.LogNormalWithPrefix(logger.NORMAL, "Task Linked: ", api)
		// local registry
		c.Publisher.Handle(router, api, tskFunc)
	}
	router.HandleFunc("/", handlerFuncBookKeeperEntry)
	return router
}

// current RPC port
func Default() *Agent {
	return &Agent{"localhost", GetPort(), true, ""}
}

// local is the local config of server
type ContactBook struct {
	Agents                 []Agent           `json:"agents,omitempty"`
	Local                  Agent             `json:"local,omitempty"`
	Coordinator            Agent             `json:"coordinator,omitempty"`
	ForbidPointToPointConn bool              `json:"forbidPointToPointConn,omitempty"`
	IsCoordinator          bool              `json:"isCoordinator,omitempty"`
	TimeStamp              int64             `json:"timestamp,omitempty"`
	Publisher              *server.Publisher `json:"-"`
	Logger                 *logger.Logger    `json:"-"`
}

// agent is the network config of server
type Agent struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Alive bool   `json:"alive"`
	API   string `json:"api,omitempty"`
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

	logger.LogNormal("Successfully Load Config File...")
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
			c.Agents[idx] = Agent{c.Agents[idx].IP, c.Agents[idx].Port, true, ""}
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
			c.Agents[index] = Agent{c.Agents[index].IP, c.Agents[index].Port, false, ""}
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
				logger.LogWarning("Connection Failed While Bridging...")
				continue
			}
			var u ContactBook
			err = client.Signal(c, &u)
			if err != nil {
				logger.LogWarning("Calling Method Failed While Bridging...")
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
				logger.LogWarning("Connection Failed While Disconnecting...")
				continue
			}
			var u ContactBook
			err = client.Disconnect(c, &u)
			if err != nil {
				logger.LogWarning("Calling Method Failed While Disconnecting...")
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
				logger.LogWarning("Connection Failed While Terminating...")
				continue
			}
			var u ContactBook
			err = client.Terminate(c, &u)
			if err != nil {
				logger.LogWarning("Calling Method Failed While Terminating...")
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

func (c *ContactBook) LaunchServer() {
	go func() {
		methods := new(LocalMethods)
		server := rpc.NewServer()
		RegisterRemote(server, methods)
		// debug setup for regcenter to check services alive
		server.HandleHTTP("/", "debug")
		l, e := net.Listen("tcp", ":"+strconv.Itoa(c.Local.Port))
		if e != nil {
			logger.LogErrorWithPrefix(logger.NORMAL, "Listen Error:", e)
		}
		http.Serve(l, nil)
	}()
	return
}

func LaunchClient(endpoint string, port int) (*RPCClient, error) {
	clientContact := Agent{endpoint, port, true, ""}
	client, err := rpc.DialHTTP("tcp", clientContact.GetFullIP())
	if err != nil {
		logger.LogErrorWithPrefix(logger.NORMAL, "Dialing:", err)
		return &RPCClient{}, err
	}
	agent := &RPCClient{Client: client}
	return agent, nil
}

func (c *Agent) GetFullIP() string {
	return c.IP + ":" + strconv.Itoa(c.Port)
}

func (c *Agent) GetFullEndPoint() string {
	return c.IP + ":" + strconv.Itoa(c.Port) + "/" + c.API
}

func (c *Agent) IsEqualTo(another *Agent) bool {
	return c.GetFullIP() == another.GetFullIP()
}

func UnmarshalAgents(original []interface{}) []Agent {
	var agents []Agent
	for _, o := range original {
		oo := o.(map[string]interface{})
		if oo["api"] != nil {
			agents = append(agents, Agent{oo["ip"].(string), int(oo["port"].(float64)), oo["alive"].(bool), oo["api"].(string)})
			continue
		}
		agents = append(agents, Agent{oo["ip"].(string), int(oo["port"].(float64)), oo["alive"].(bool), ""})
	}
	return agents
}

func handlerFuncBookKeeperEntry(w http.ResponseWriter, r *http.Request) {
	logger.LogNormal("Bookkeeper Println...")
}
