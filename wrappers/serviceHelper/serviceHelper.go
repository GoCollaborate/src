package serviceHelper

import (
	// "bytes"
	// "encoding/json"
	// "github.com/GoCollaborate/src/artifacts/card"
	// "github.com/GoCollaborate/src/artifacts/restful"
	"github.com/GoCollaborate/src/artifacts/service"
	// "github.com/GoCollaborate/src/constants"
	// "github.com/GoCollaborate/src/utils"
	// "io"
	// "net/http"
)

func ModeInterpret(original interface{}) service.Mode {
	if original == nil {
		return service.ClbtModeNormal
	}
	var m service.Mode
	omode := original.(string)
	switch omode {
	case "ClbtModeOnlyRegister":
		m = service.ClbtModeOnlyRegister
	case "ClbtModeOnlySubscribe":
		m = service.ClbtModeOnlySubscribe
	case "LBModeTokenHash":
		m = service.LBModeTokenHash
	case "LBModeRandom":
		m = service.LBModeRandom
	case "LBModeLeastActive":
		m = service.LBModeLeastActive
	case "LBModeRoundRobin":
		m = service.LBModeRoundRobin
	default:
		m = service.ClbtModeNormal
	}
	return m
}
