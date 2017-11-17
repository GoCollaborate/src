package serviceHelper

import (
	"bytes"
	"encoding/json"
	"github.com/GoCollaborate/src/artifacts/card"
	"github.com/GoCollaborate/src/artifacts/restful"
	"github.com/GoCollaborate/src/artifacts/service"
	"github.com/GoCollaborate/src/constants"
	"github.com/GoCollaborate/src/utils"
	"io"
	"net/http"
)

func ModeInterpret(original interface{}) service.Mode {
	if original == nil {
		return service.ClbtModeNormal
	}
	var m service.Mode
	omode := original.(string)
	switch omode {
	case "ClbtModeOnlyRegister":
		m = service.ClbtModeOnlyRegister
	case "ClbtModeOnlySubscribe":
		m = service.ClbtModeOnlySubscribe
	case "LBModeTokenHash":
		m = service.LBModeTokenHash
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

func MarshalServiceResourceToByteStream(srvs map[string]*service.Service) ([]byte, error) {
	payload := service.ServicePayload{
		Data: MarshalServiceToStandardResource(srvs),
	}
	return json.Marshal(payload)
}

func MarshalServiceResourceToByteStreamReader(srvs map[string]*service.Service) (io.Reader, error) {
	b, err := MarshalServiceResourceToByteStream(srvs)
	return bytes.NewReader(b), err
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

func MarshalHeartbeatToStandardResource(card *card.Card, services map[string]*service.Service) []service.HeartbeatResource {
	var resources []service.HeartbeatResource
	for srvID, _ := range services {
		resources = append(resources, service.HeartbeatResource{
			Resource: &restful.Resource{
				ID:            srvID,
				Type:          "heartbeat",
				Relationships: []restful.Relationship{}},
			Attributes: service.Heartbeat{
				Agent: *card,
			}})
	}
	return resources
}

func MarshalHeartbeatToByteStream(card *card.Card, services map[string]*service.Service) ([]byte, error) {
	payload := service.HeartbeatPayload{
		Data: MarshalHeartbeatToStandardResource(card, services),
	}
	return json.Marshal(payload)
}

func MarshalHeartbeatToByteStreamReader(card *card.Card, services map[string]*service.Service) (io.Reader, error) {
	b, err := MarshalHeartbeatToByteStream(card, services)
	return bytes.NewReader(b), err
}

func DecodeService(bytes []byte) *service.ServicePayload {
	payload := service.ServicePayload{}
	json.Unmarshal(bytes, &payload)
	return &payload
}

func DecodeRegistry(bytes []byte) *service.RegistryPayload {
	payload := service.RegistryPayload{}
	json.Unmarshal(bytes, &payload)
	return &payload
}

func DecodeSubscription(bytes []byte) *service.SubscriptionPayload {
	payload := service.SubscriptionPayload{}
	json.Unmarshal(bytes, &payload)
	return &payload
}

func DecodeHeartbeat(bytes []byte) *service.HeartbeatPayload {
	payload := service.HeartbeatPayload{}
	json.Unmarshal(bytes, &payload)
	return &payload
}

func SendServiceWith(w http.ResponseWriter, srvPayload *service.ServicePayload, status int) error {
	mal, err := json.Marshal(*srvPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.HeaderContentTypeJSON)
	utils.AdaptHTTPWithStatus(w, status)
	io.WriteString(w, string(mal))
	return nil
}

func SendRegistryWith(w http.ResponseWriter, regPayload *service.RegistryPayload, status int) error {
	mal, err := json.Marshal(*regPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.HeaderContentTypeJSON)
	utils.AdaptHTTPWithStatus(w, status)
	io.WriteString(w, string(mal))
	return nil
}

func SendSubscriptionWith(w http.ResponseWriter, scpPayload *service.SubscriptionPayload, status int) error {
	mal, err := json.Marshal(*scpPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.HeaderContentTypeJSON)
	utils.AdaptHTTPWithStatus(w, status)
	io.WriteString(w, string(mal))
	return nil
}

func SendQueryWith(w http.ResponseWriter, queryPayload *service.QueryPayload, status int) error {
	mal, err := json.Marshal(*queryPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.HeaderContentTypeJSON)
	utils.AdaptHTTPWithStatus(w, status)
	io.WriteString(w, string(mal))
	return nil
}

func SendHeartbeatWith(w http.ResponseWriter, status int) error {
	utils.AdaptHTTPWithHeader(w, constants.HeaderContentTypeJSON)
	utils.AdaptHTTPWithStatus(w, status)
	io.WriteString(w, "")
	return nil
}
