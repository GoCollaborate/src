package service

import (
	"github.com/GoCollaborate/src/artifacts/restful"
)

type HeartbeatResource struct {
	*restful.Resource
	Attributes Heartbeat `json:"attributes"`
}

type HeartbeatPayload struct {
	Data     []HeartbeatResource `json:"data"`
	Included []HeartbeatResource `json:"included,omitempty"`
	Links    restful.Links       `json:"links,omitempty"`
}
