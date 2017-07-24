package remote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

const configPath = "./config.json"

// current RPC port
func Default() *Agent {
	return &Agent{"localhost", GetPort(), true}
}

// local is the local config of server
type Config struct {
	Agents    []Agent `json:"agents,omitempty"`
	TimeStamp int64   `json:"timestamp,omitempty"`
	Local     Agent   `json:"local,omitempty"`
}

// agent is the network config of server
type Agent struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Alive bool   `json:"alive"`
}

func Populate(cfg *Config) error {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	// unmarshal, overwrite default if already existed in config file
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *Config) Load() (*Config, error) {
	err := Populate(c)
	if err != nil {
		return c, err
	}

	fmt.Println("Localaddress:")
	ip := GetLocalIP()
	fmt.Println(ip)

	var exist bool = false
	var idx int

	fmt.Println("Successfully Load Config File...")
	fmt.Println("Hosts:")

	for i, h := range c.Agents {
		if h.IsEqualTo(&c.Local) {
			exist = true
			idx = i
		}
		fmt.Print(h.IP)
		fmt.Print(":")
		fmt.Println(h.Port)
	}

	// update if not exist
	if !exist {
		c.Agents = append(c.Agents, c.Local)
		c.Sync()
	} else {
		// activate if not alive
		if !c.Agents[idx].Alive {
			c.Agents[idx] = Agent{c.Agents[idx].IP, c.Agents[idx].Port, true}
			c.Sync()
		}
	}
	return c, nil
}

func (c *Config) RemoteLoad() (*Config, error) {
	var localConfig *Config = new(Config)
	err := Populate(localConfig)
	if err != nil {
		return c, err
	}
	var exist bool = false
	var update bool = false

	update = Compare(localConfig, c)

	for _, h := range c.Agents {
		if h.IsEqualTo(&localConfig.Local) {
			exist = true
		}
	}
	// update if not exist
	if !exist {
		c.Agents = append(c.Agents, c.Local)
		localConfig.Agents = append(c.Agents, c.Local)
		update = true
	}
	if update {
		localConfig.Sync()
	}
	return c, nil
}

func (c *Config) RemoteDisconnect() (*Config, error) {
	var localConfig *Config = new(Config)
	err := Populate(localConfig)
	if err != nil {
		return c, err
	}
	var index = -1
	var update bool = false

	update = Compare(localConfig, c)

	for i, h := range c.Agents {
		if h.IsEqualTo(&c.Local) {
			index = i
		}
	}
	// update if not exist
	if index < 0 {
		c.Agents = append(c.Agents, c.Local)
		localConfig.Agents = append(c.Agents, c.Local)
		update = true
	} else {
		// deactivate if alive
		if c.Agents[index].Alive {
			c.Agents[index] = Agent{c.Agents[index].IP, c.Agents[index].Port, false}
			c.Stamp()
			for _, a := range localConfig.Agents {
				if a.IsEqualTo(&c.Agents[index]) {
					a.Alive = false
				}
			}
			update = true
		}
	}
	if update {
		localConfig.Sync()
	}
	return c, nil
}

func (c *Config) RemoteTerminate() (*Config, error) {
	var localConfig *Config = new(Config)
	err := Populate(localConfig)
	if err != nil {
		return c, err
	}
	var index int = -1
	var update bool = false

	for i, h := range localConfig.Agents {
		if h.IsEqualTo(&c.Local) {
			index = i
			update = true
		}
	}
	// update if exist
	if index >= 0 {
		localConfig.RemoveAgent(index)
		index = -1
	}
	for j, h := range c.Agents {
		if h.IsEqualTo(&c.Local) {
			index = j
		}
	}
	fmt.Println(c)
	if index >= 0 {
		c.RemoveAgent(index)
		c.Stamp()
	}

	fmt.Println(c)

	if update {
		localConfig.Sync()
	}
	return c, nil
}

func (c *Config) Bridge() {
	for _, e := range c.Agents {
		if e.Alive && (!e.IsEqualTo(&c.Local)) {
			client, err := LaunchClient(e.IP, e.Port)
			if err != nil {
				fmt.Println("Connection Failed While Bridging...")
				continue
			}
			var u Config
			err = client.Signal(c, &u)
			if err != nil {
				fmt.Println("Calling Method Failed While Bridging...")
				continue
			}
			u.Update(c).WriteStream()
		}
	}
}

func (c *Config) Disconnect() {
	for _, e := range c.Agents {
		if e.Alive && (!e.IsEqualTo(&c.Local)) {
			client, err := LaunchClient(e.IP, e.Port)
			if err != nil {
				fmt.Println("Connection Failed While Disconnecting...")
				continue
			}
			var u Config
			err = client.Disconnect(c, &u)
			if err != nil {
				fmt.Println("Calling Method Failed While Disconnecting...")
				continue
			}
			u.Update(c).WriteStream()
		}
	}
}

func (c *Config) Terminate() {
	for _, e := range c.Agents {
		if e.Alive && (!e.IsEqualTo(&c.Local)) {
			client, err := LaunchClient(e.IP, e.Port)
			if err != nil {
				fmt.Println("Connection Failed While Terminating...")
				continue
			}
			var u Config
			err = client.Terminate(c, &u)
			if err != nil {
				fmt.Println("Calling Method Failed While Terminating...")
				continue
			}
			u.Update(c).WriteStream()
		}
	}
}

// Sync will also update TimeStamp
func (c *Config) Sync() {
	c.TimeStamp = time.Now().Unix()
	mal, err := json.Marshal(&c)
	err = ioutil.WriteFile(configPath, mal, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
}

// Update() does not update TimeStamp
func (c *Config) Update(update *Config) (cfg *Config) {
	Compare(update, c)
	return c
}

func (c *Config) WriteStream() {
	mal, err := json.Marshal(&c)
	err = ioutil.WriteFile(configPath, mal, os.ModeExclusive)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Config) RemoveAgent(index int) *Config {
	c.Agents = append(c.Agents[:(index)], c.Agents[(index+1):]...)
	return c
}

func (c *Config) Stamp() *Config {
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

func Compare(a *Config, b *Config) bool {
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

func (c *Config) LaunchServer() {
	go func() {
		methods := new(LocalMethods)
		server := rpc.NewServer()
		RegisterRemote(server, methods)
		server.HandleHTTP("/", "/debug")
		l, e := net.Listen("tcp", ":"+strconv.Itoa(c.Local.Port))
		if e != nil {
			fmt.Println("Listen Error:", e)
		}
		http.Serve(l, nil)
	}()
}

func LaunchClient(endpoint string, port int) (*RPCClient, error) {
	clientConfig := Agent{endpoint, port, true}
	client, err := rpc.DialHTTP("tcp", clientConfig.GetFullIP())
	if err != nil {
		fmt.Println("Dialing:", err)
		return &RPCClient{}, err
	}
	agent := &RPCClient{Client: client}
	return agent, nil
}

func (c *Agent) GetFullIP() string {
	return c.IP + ":" + strconv.Itoa(c.Port)
}

func (c *Agent) IsEqualTo(another *Agent) bool {
	return c.GetFullIP() == another.GetFullIP()
}
