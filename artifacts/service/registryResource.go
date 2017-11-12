package service

import (
	"github.com/GoCollaborate/src/artifacts/restful"
)

type RegistryResource struct {
	*restful.Resource
	Attributes Registry `json:"attributes"`
}

type RegistryPayload struct {
	Data     []RegistryResource `json:"data"`
	Included []RegistryResource `json:"included,omitempty"`
	Links    restful.Links      `json:"links,omitempty"`
}
