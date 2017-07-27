package remote

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
