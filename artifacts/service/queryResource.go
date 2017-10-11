package service

import (
	"github.com/GoCollaborate/artifacts/restful"
)

type QueryResource struct {
	Resource   *restful.Resource
	Attributes Query `json:"attributes"`
}

type QueryPayload struct {
	Data     []QueryResource `json:"data"`
	Included []QueryResource `json:"included,omitempty"`
	Links    restful.Links   `json:"links,omitempty"`
}
