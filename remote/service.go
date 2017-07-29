package remote

import (
	"github.com/GoCollaborate/constants"
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
	ServiceID        string      `json:"serviceid"`
	Description      string      `json:"description"`
	Parameters       []Parameter `json:"parameters"`
	RegList          []Agent     `json:"registers"`
	SbscrbList       []string    `json:"subscribers"`  // subscriber tokens
	Dependencies     []string    `json:"dependencies"` // dependent ServiceIDs
	Mode             mode        `json:"mode"`
	LoadBalanceMode  mode        `json:"load_balance_mode"`
	Version          string      `json:"version"`
	PlatformVersion  string      `json:"platform_version"`
	LastAssignedTo   string      `json:"last_assigned_to"`
	LastAssignedTime int64       `json:"last_assigned_time"`
}

type Mode interface {
	GetStatus()
}

type mode int

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

func (s *Service) GetDependencies() []string {
	return s.Dependencies
}

func (s *Service) SetLoadBalanceMode(m *mode) mode {
	s.LoadBalanceMode = *m
	return *m
}

func (s *Service) GetLoadBalanceMode() mode {
	return s.LoadBalanceMode
}

func (s *Service) Register(agt *Agent) error {
	for _, r := range s.RegList {
		if agt.IsEqualTo(&r) {
			return constants.ErrConflictRegister
		}
	}
	s.RegList = append(s.RegList, *agt)
	return nil
}

// de-register function will not preserve the original order of registers
func (s *Service) DeRegister(agt *Agent) error {
	y := s.RegList[:0]
	for i, x := range s.RegList {
		if agt.IsEqualTo(&x) {
			y = append(s.RegList[:i], s.RegList[i+1:]...)
			s.RegList = y
			return nil
		}
	}
	return constants.ErrNoSubscriber
}

func (s *Service) Subscribe(token string) error {
	for _, sbscr := range s.SbscrbList {
		if token == sbscr {
			return constants.ErrConflictSubscriber
		}
	}
	s.SbscrbList = append(s.SbscrbList, token)
	return nil
}

func (s *Service) UnSubscribe(token string) error {
	y := s.SbscrbList[:0]
	for i, x := range s.SbscrbList {
		if x == token {
			y = append(s.SbscrbList[:i], s.SbscrbList[i+1:]...)
			s.SbscrbList = y
			return nil
		}
	}
	return constants.ErrNoSubscriber
}
