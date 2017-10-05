package serviceHelper

import (
	"github.com/GoCollaborate/artifacts/service"
)

func UnmarshalMode(original interface{}) service.Mode {
	if original == nil {
		return service.RPCServerModeNormal
	}
	var m service.Mode
	omode := original.(string)
	switch omode {
	case "RPCServerModeOnlyRegister":
		m = service.RPCServerModeOnlyRegister
	case "RPCServerModeStatic":
		m = service.RPCServerModeStatic
	case "RPCServerModeRandomLoadBalance":
		m = service.RPCServerModeRandomLoadBalance
	case "RPCServerModeLeastActiveLoadBalance":
		m = service.RPCServerModeLeastActiveLoadBalance
	case "RPCClientModeOnlySubscribe":
		m = service.RPCClientModeOnlySubscribe
	case "RPCClientModePointToPoint":
		m = service.RPCClientModePointToPoint
	default:
		m = service.RPCServerModeNormal
	}
	return m
}
