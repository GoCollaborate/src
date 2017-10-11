package serviceHelper

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/card"
	"github.com/GoCollaborate/artifacts/restful"
	"github.com/GoCollaborate/artifacts/service"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/utils"
	"io"
	"net/http"
)

func UnmarshalMode(original interface{}) service.Mode {
	if original == nil {
		return service.ClbtModeNormal
	}
	var m service.Mode
	omode := original.(string)
	switch omode {
	case "LBModeIPHash":
		m = service.LBModeIPHash
	case "ClbtModeOnlyRegister":
		m = service.ClbtModeOnlyRegister
	case "ClbtModeOnlySubscribe":
		m = service.ClbtModeOnlySubscribe
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

func MarshalServiceToStandardResource(srvs map[string]*service.Service) []service.ServiceResource {
	var resources []service.ServiceResource
	for _, srv := range srvs {
		resources = append(resources, service.ServiceResource{
			Resource: &restful.Resource{
				ID:            srv.ServiceID,
				Type:          "service",
				Relationships: []restful.Relationship{}},
			Attributes: *srv})
	}
	return resources
}

func MarshalCardToStandardResource(srvID string, card *card.Card) []service.QueryResource {
	var resources []service.QueryResource
	resources = append(resources, service.QueryResource{
		Resource: &restful.Resource{
			ID:            srvID,
			Type:          "query",
			Relationships: []restful.Relationship{}},
		Attributes: service.Query{
			Agent: *card,
		}})
	return resources
}

func DecodeService(js string) *service.ServicePayload {
	payload := service.ServicePayload{}
	json.Unmarshal([]byte(js), &payload)
	return &payload
}

func DecodeRegistry(js string) *service.RegistryPayload {
	payload := service.RegistryPayload{}
	json.Unmarshal([]byte(js), &payload)
	return &payload
}

func DecodeSubscription(js string) *service.SubscriptionPayload {
	payload := service.SubscriptionPayload{}
	json.Unmarshal([]byte(js), &payload)
	return &payload
}

func SendServiceWith(w http.ResponseWriter, srvPayload *service.ServicePayload, header constants.Header) error {
	mal, err := json.Marshal(*srvPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, header)
	io.WriteString(w, string(mal))
	return nil
}

func SendRegistryWith(w http.ResponseWriter, regPayload *service.RegistryPayload, header constants.Header) error {
	mal, err := json.Marshal(*regPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, header)
	io.WriteString(w, string(mal))
	return nil
}

func SendSubscriptionWith(w http.ResponseWriter, scpPayload *service.SubscriptionPayload, header constants.Header) error {
	mal, err := json.Marshal(*scpPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, header)
	io.WriteString(w, string(mal))
	return nil
}

func SendQueryWith(w http.ResponseWriter, queryPayload *service.QueryPayload, header constants.Header) error {
	mal, err := json.Marshal(*queryPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, header)
	io.WriteString(w, string(mal))
	return nil
}
