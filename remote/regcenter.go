package remote

import (
	"runtime"
	"time"
)

// remote settings
const (
	RPCServerModeNormal                 mode = iota
	RPCServerModeRegCenter                   // turn RPC server of this mode to RegCenter
	RPCServerModeOnlyRegister                // collaborator only registers service, service is not accessible until it has been changed to RPCServerModeNormal
	RPCServerModeStatic                      // coordinating agent will not automatically manage this server
	RPCServerModeRandomLoadBalance           // assign tasks as per weighted probability
	RPCServerModeLeastActiveLoadBalance      // assign tasks to active responders
	RPCClientModeOnlySubscribe               // client only subscribes to collaborator service at RegCenter, yet still collaborates with direct point to point connection
	RPCClientModePointToPoint                // this will provide point to point connection between client and collaborator
)

type Service struct {
	Name             string
	Arguments        []Argument
	RegList          []Agent
	SbscrbList       []string // subscriber tokens
	Mode             mode
	LoadBalanceMode  mode
	LastAssignedTo   Agent
	LastAssignedTime int64
}

type RegCenter struct {
	Created  int64 // time in epoch milliseconds indicating when the registry center was created
	Modified int64 // time in epoch milliseconds indicating when the registry center was last modified
	Services map[string]Service
	Port     int
	GoVer    string
}

type Mode interface {
	GetStatus()
}

type mode int

func NewRegCenter(port int) *RegCenter {
	return &RegCenter{time.Now().Unix(), time.Now().Unix(), map[string]Service{}, port, runtime.Version()}
}

func (m *mode) GetMode() mode {
	return *m
}

func (s *Service) SetMode(m *mode) mode {
	s.Mode = *m
	return *m
}

func (s *Service) GetMode() mode {
	return s.Mode
}

func (s *Service) SetLoadBalanceMode(m *mode) mode {
	s.LoadBalanceMode = *m
	return *m
}

func (s *Service) GetLoadBalanceMode() mode {
	return s.LoadBalanceMode
}

func (rc *RegCenter) Listen() {

}

func (rc *RegCenter) HeartBeatAll() {
	// send heartbeat signals to all service providers
	// return which of them have already expired
}

func (rc *RegCenter) RefreshList() {
	rc.HeartBeatAll()
	// acquire logs from heartbeat detection and remove those which have already expired
}

func (rc *RegCenter) Dump() {
	// dump registry data to local dat file
	// constants.DataStorePath
}
