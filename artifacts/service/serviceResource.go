package service

import (
	"github.com/GoCollaborate/artifacts/restful"
)

type ServiceResource struct {
	Resource   *restful.Resource
	Attributes Service `json:"attributes"`
}

type ServicePayload struct {
	Data     []ServiceResource `json:"data"`
	Included []ServiceResource `json:"included,omitempty"`
	Links    restful.Links     `json:"links,omitempty"`
}
