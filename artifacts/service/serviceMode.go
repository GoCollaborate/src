package service

// remote settings
const (
	RPCServerModeNormal                 Mode = iota
	RPCServerModeOnlyRegister                // collaborator only registers service, service is not accessible until it has been changed to RPCServerModeNormal
	RPCServerModeStatic                      // coordinating Card will not automatically manage this server
	RPCServerModeRandomLoadBalance           // assign tasks as per weighted probability
	RPCServerModeLeastActiveLoadBalance      // assign tasks to active responders
	RPCClientModeOnlySubscribe               // client only subscribes to collaborator service at RegCenter, yet still collaborates with direct point to point connection
	RPCClientModePointToPoint                // this will provide point to point connection between client and collaborator
)

type Mode int

func (m *Mode) GetMode() Mode {
	return *m
}
