package service

import (
	. "github.com/GoCollaborate/src/artifacts/resources"
)

type ServiceResource struct {
	*Resource
	Attributes *Service `json:"attributes"`
}

func NewServiceResource(service *Service) *ServiceResource {
	return &ServiceResource{
		&Resource{
			Id:            service.ServiceID,
			Type:          "service",
			Relationships: map[string]*Relationship{},
		},
		service,
	}
}

type ServicesResource struct {
	*Resource
	Attributes *map[string]*Service `json:"attributes"`
}

func NewServicesResource(services *map[string]*Service) *ServicesResource {
	return &ServicesResource{
		&Resource{
			Type:          "services",
			Relationships: map[string]*Relationship{},
		},
		services,
	}
}
