package service

import (
	"github.com/GoCollaborate/src/artifacts/restful"
)

type ServiceResource struct {
	*restful.Resource
	Attributes Service `json:"attributes"`
}

type ServicePayload struct {
	Data     []ServiceResource `json:"data"`
	Included []ServiceResource `json:"included,omitempty"`
	Links    restful.Links     `json:"links,omitempty"`
}
