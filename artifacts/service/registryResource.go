package service

import (
	"github.com/GoCollaborate/artifacts/restful"
)

type RegistryResource struct {
	Resource   *restful.Resource
	Attributes Registry `json:"attributes"`
}

type RegistryPayload struct {
	Data     []RegistryResource `json:"data"`
	Included []RegistryResource `json:"included,omitempty"`
	Links    restful.Links      `json:"links,omitempty"`
}
