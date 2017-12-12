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

func NewServiceResources(services ...*Service) []*ServiceResource {
	var (
		resources = []*ServiceResource{}
	)

	for _, s := range services {
		resources = append(resources, &ServiceResource{
			&Resource{
				Id:            s.ServiceID,
				Type:          "service",
				Relationships: map[string]*Relationship{},
			},
			s,
		})
	}

	return resources
}

func NewServiceResourcesFromMap(services *map[string]*Service) []*ServiceResource {
	var (
		resources = []*ServiceResource{}
	)

	for _, s := range *services {
		resources = append(resources, &ServiceResource{
			&Resource{
				Id:            s.ServiceID,
				Type:          "service",
				Relationships: map[string]*Relationship{},
			},
			s,
		})
	}

	return resources
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
