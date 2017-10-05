package restfulHelper

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/card"
	"github.com/GoCollaborate/artifacts/restful"
	"github.com/GoCollaborate/artifacts/service"
)

func Decode(js string) *restful.ExposurePayload {
	payload := restful.ExposurePayload{}
	json.Unmarshal([]byte(js), &payload)
	return &payload
}

func MarshalServiceToStandardResource(srvs map[string]*service.Service) []restful.Resource {
	var resources []restful.Resource
	for _, srv := range srvs {
		resources = append(resources, restful.Resource{
			"service",
			srv.ServiceID,
			map[string]interface{}{
				"description":       srv.Description,
				"parameters":        srv.Parameters,
				"registers":         srv.RegList,
				"subscribers":       srv.SbscrbList,
				"dependencies":      srv.Dependencies,
				"mode":              srv.Mode,
				"load_balance_mode": srv.LoadBalanceMode,
				"version":           srv.Version,
				"platform_version":  srv.PlatformVersion,
			}, []restful.Relationship{}})
	}
	return resources
}

func MarshalCardToStandardResource(srvID string, agt *card.Card) []restful.Resource {
	return []restful.Resource{restful.Resource{"query", srvID, map[string]interface{}{
		"enpoint": agt.GetFullEndPoint(), "ip": agt.IP, "port": agt.Port}, []restful.Relationship{}}}
}
