package coordinator

import (
	"encoding/json"
	"github.com/GoCollaborate/src/artifacts/card"
	"github.com/GoCollaborate/src/artifacts/restful"
	"github.com/GoCollaborate/src/artifacts/service"
	"github.com/GoCollaborate/src/constants"
	"github.com/GoCollaborate/src/logger"
	"github.com/GoCollaborate/src/utils"
	"github.com/GoCollaborate/src/wrappers/restfulHelper"
	"github.com/GoCollaborate/src/wrappers/serviceHelper"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	TokenLength   = 10
	ServiceLength = 15
)

type Coordinator struct {
	Created  int64                          `json:"created"`  // time in epoch milliseconds indicating when the registry center was created
	Modified int64                          `json:"modified"` // time in epoch milliseconds indicating when the registry center was last modified
	Services map[string]*service.Service    `json:"services"` // a map of ServiceID to services
	Clusters map[string]map[string]struct{} `json:"clusters"` // a set of cluster services
	Port     int                            `json:"port"`
	GoVer    string                         `json:"gover"`
}

func (cdnt *Coordinator) Handle(router *mux.Router) *mux.Router {

	router.HandleFunc("/{version}/services", HandlerFuncGetServices).Methods("GET")

	// client launch request to call a service
	router.HandleFunc("/{version}/query/{srvid}/{token}", HandlerFuncQueryService).Methods("GET")

	router.HandleFunc("/{version}/services/{srvid}", HandlerFuncGetService).Methods("GET")
	router.HandleFunc("/{version}/services/{srvid}", HandlerFuncDeleteService).Methods("DELETE")

	router.HandleFunc("/{version}/services/{srvid}/registry", HandlerFuncRegisterService).Methods("POST")
	router.HandleFunc("/{version}/services/{srvid}/subscription", HandlerFuncSubscribeService).Methods("POST")

	// single deletion [DELETE BODY is not explicitly specified in HTTP 1.1,
	// so it's better not to use its message-body]
	router.HandleFunc("/{version}/services/{srvid}/registry/single/{ip}/{port}", HandlerFuncDeRegisterService).Methods("DELETE")
	router.HandleFunc("/{version}/services/{srvid}/subscription/single/{token}", HandlerFuncUnSubscribeService).Methods("DELETE")

	// bulk deletion
	router.HandleFunc("/{version}/services/{srvid}/registry", HandlerFuncDeRegisterServiceAll).Methods("DELETE")
	router.HandleFunc("/{version}/services/{srvid}/subscription", HandlerFuncUnSubscribeServiceAll).Methods("DELETE")

	// create services
	router.HandleFunc("/{version}/services", HandlerFuncPostServices).Methods("POST")
	// alter services
	router.HandleFunc("/{version}/services", HandlerFuncAlterServices).Methods("PUT")

	// router.HandleFunc("/services/{srvid}/heartbeat", HandlerFuncHeartbeatServices).Methods("GET")

	// clustering apis
	router.HandleFunc("/{version}/cluster/{cid}/services", HanlderFuncGetClusterServices).Methods("GET")
	// include service in cluster
	// no message body is required
	router.HandleFunc("/{version}/cluster/{cid}/services", HanlderFuncPostClusterServices).Methods("POST")
	// services should be decoupled from cluster
	router.HandleFunc("/{version}/services/heartbeat", HandlerFuncHeartbeatServices).Methods("POST")

	go func() {
		<-time.After(constants.DefaultCollaboratorExpiryInterval)
		cdnt.Clean()
		cdnt.Dump()
	}()
	return router
}

func (cdnt *Coordinator) GetService(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := cdnt.Services[srvID]; ok {
		services := map[string]*service.Service{}
		services[srvID] = cdnt.Services[srvID]
		payload := service.ServicePayload{
			Data:     serviceHelper.MarshalServiceToStandardResource(services),
			Included: []service.ServiceResource{},
			Links: restful.Links{
				Self: r.URL.String(),
			},
		}
		return serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
	}
	return restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
}

func (cdnt *Coordinator) GetServices(w http.ResponseWriter, r *http.Request) error {
	// look up service list, returning service providers
	payload := service.ServicePayload{
		Data:     serviceHelper.MarshalServiceToStandardResource(cdnt.Services),
		Included: []service.ServiceResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}
	return serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
}

// POST
func (cdnt *Coordinator) CreateService(w http.ResponseWriter, bytes []byte) error {
	// parse request json and create service here
	payload := serviceHelper.DecodeService(bytes)
	var resData []service.ServiceResource
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			id := utils.RandStringBytesMaskImprSrc(ServiceLength)
			svr := service.NewServiceFrom(&dat.Attributes)
			svr.ServiceID = id
			cdnt.Services[id] = svr
			dat.ID = id
			dat.Attributes = *svr
			resData = append(resData, dat)
		}
	}
	payload.Data = resData
	serviceHelper.SendServiceWith(w, payload, http.StatusCreated)
	return nil
}

// POST
func (cdnt *Coordinator) RegisterService(w http.ResponseWriter, bytes []byte, srvID string) error {
	payload := serviceHelper.DecodeRegistry(bytes)
	for _, dat := range payload.Data {
		if dat.Type == "registry" {
			Cards := dat.Attributes.Cards
			if _, ok := cdnt.Services[srvID]; !ok {
				restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
				return constants.ErrNoService
			}
			for _, Card := range Cards {
				err := cdnt.Services[srvID].Register(&Card)
				if err != nil {
					restfulHelper.SendErrorWith(w, restful.Error409Conflict(), http.StatusConflict)
					return err
				}
			}
		}
	}
	serviceHelper.SendRegistryWith(w, payload, http.StatusOK)
	return nil
}

// POST
func (cdnt *Coordinator) SubscribeService(w http.ResponseWriter, bytes []byte, srvID string) error {
	var res []service.SubscriptionResource
	payload := serviceHelper.DecodeSubscription(bytes)
	for _, dat := range payload.Data {
		if dat.Type == "subscription" {
			if _, ok := cdnt.Services[srvID]; ok {
				tk := utils.RandStringBytesMaskImprSrc(TokenLength)
				err := cdnt.Services[srvID].Subscribe(tk)
				if err != nil {
					restfulHelper.SendErrorWith(w, restful.Error409Conflict(), http.StatusConflict)
					return err
				}
				dat.Attributes.Token = tk
				res = append(res, dat)
			}
		}
	}
	payload.Data = res
	serviceHelper.SendSubscriptionWith(w, payload, http.StatusOK)
	return nil
}

// DELETE
func (cdnt *Coordinator) DeleteService(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := cdnt.Services[srvID]; ok {
		services := map[string]*service.Service{}
		services[srvID] = cdnt.Services[srvID]
		payload := service.ServicePayload{
			Data:     serviceHelper.MarshalServiceToStandardResource(services),
			Included: []service.ServiceResource{},
			Links: restful.Links{
				Self: r.URL.String(),
			},
		}
		delete(cdnt.Services, srvID)
		serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
		return nil
	}
	restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
	return constants.ErrNoService

}

// DELETE
func (cdnt *Coordinator) DeRegisterService(w http.ResponseWriter, r *http.Request, srvID string, cd *card.Card) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		return constants.ErrNoService
	}

	err := cdnt.Services[srvID].DeRegister(cd)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		return err
	}

	services := map[string]*service.Service{}
	services[srvID] = cdnt.Services[srvID]

	payload := service.ServicePayload{
		Data:     serviceHelper.MarshalServiceToStandardResource(services),
		Included: []service.ServiceResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}

	serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
	return nil
}

// DELETE
func (cdnt *Coordinator) DeRegisterServiceAll(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		return constants.ErrNoService
	}
	services := map[string]*service.Service{}
	services[srvID] = cdnt.Services[srvID]
	payload := service.ServicePayload{
		Data:     serviceHelper.MarshalServiceToStandardResource(services),
		Included: []service.ServiceResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}
	cdnt.Services[srvID].DeRegisterAll()

	serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
	return nil
}

// DELETE
func (cdnt *Coordinator) UnSubscribeService(w http.ResponseWriter, r *http.Request, srvID string, token string) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		return constants.ErrNoService
	}

	err := cdnt.Services[srvID].UnSubscribe(token)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		return err
	}

	services := map[string]*service.Service{}
	services[srvID] = cdnt.Services[srvID]

	payload := service.ServicePayload{
		Data:     serviceHelper.MarshalServiceToStandardResource(services),
		Included: []service.ServiceResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}

	serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
	return nil
}

// DELETE
func (cdnt *Coordinator) UnSubscribeServiceAll(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		return constants.ErrNoService
	}
	services := map[string]*service.Service{}
	services[srvID] = cdnt.Services[srvID]
	payload := service.ServicePayload{
		Data:     serviceHelper.MarshalServiceToStandardResource(services),
		Included: []service.ServiceResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}
	cdnt.Services[srvID].UnSubscribeAll()
	serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
	return nil
}

// PUT
func (cdnt *Coordinator) AlterService(w http.ResponseWriter, bytes []byte) error {
	payload := serviceHelper.DecodeService(bytes)
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			srvID := dat.ID
			if _, ok := cdnt.Services[srvID]; !ok {
				restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
				return constants.ErrNoService
			}
			cdnt.Services[srvID] = &dat.Attributes
		}
	}
	serviceHelper.SendServiceWith(w, payload, http.StatusAccepted)
	return nil
}

func (cdnt *Coordinator) QueryService(w http.ResponseWriter, r *http.Request, srvID string, token string) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		return constants.ErrNoService
	}

	agt, err := cdnt.Services[srvID].Query(token)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		return err
	}

	payload := service.QueryPayload{
		Data:     serviceHelper.MarshalCardToStandardResource(srvID, agt),
		Included: []service.QueryResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}

	serviceHelper.SendQueryWith(w, &payload, http.StatusOK)
	return nil
}

// the coordinator active heartbeat is only available when all of the collaborators are running in debug mode
func (rc *Coordinator) HeartBeatAll() []card.Card {
	// send heartbeat signals to all service providers
	// return which of them have already expired
	var expired []card.Card
	for _, service := range rc.Services {
		for _, register := range service.RegList {
			// test via predeclared endpoint
			res, err := http.Get(register.GetFullIP() + "/debug/heartbeat")
			if err != nil || res.StatusCode != 200 {
				expired = append(expired, register)
			}
		}
	}
	return expired
}

func (cdnt *Coordinator) RefreshList() {
	// acquire logs from heartbeat detection and remove those which have already expired
	for _, service := range cdnt.Services {
		var activated []card.Card
		for _, register := range service.RegList {
			res, err := http.Get(register.GetFullIP() + "/debug/heartbeat")
			if err == nil && res.StatusCode == 200 {
				activated = append(activated, register)
			}
		}
		service.RegList = activated
	}
	return
}

func (cdnt *Coordinator) HeartbeatServices(w http.ResponseWriter, bytes []byte) error {
	payload := serviceHelper.DecodeHeartbeat(bytes)

	for _, dat := range payload.Data {
		if dat.Type == "heartbeat" {
			srvID := dat.ID
			if _, ok := cdnt.Services[srvID]; !ok {
				restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
				return constants.ErrNoService
			}
			cdnt.Services[srvID].Heartbeat(&dat.Attributes.Agent)
		}
	}
	return serviceHelper.SendHeartbeatWith(w, http.StatusOK)
}

func (cdnt *Coordinator) GetClusterServices(w http.ResponseWriter, r *http.Request, cID string) error {
	var (
		ok       = false
		list     = map[string]struct{}{}
		services = map[string]*service.Service{}
	)

	if list, ok = cdnt.Clusters[cID]; !ok {
		return restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
	}

	for key, _ := range list {
		services[key] = cdnt.Services[key]
	}

	payload := service.ServicePayload{
		Data:     serviceHelper.MarshalServiceToStandardResource(services),
		Included: []service.ServiceResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}
	return serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
}

func (cdnt *Coordinator) PostClusterServices(w http.ResponseWriter, r *http.Request, cID string, bytes []byte) error {
	var (
		ok       = false
		list     = map[string]struct{}{}
		services = map[string]*service.Service{}
		exists   = struct{}{}
	)

	if list, ok = cdnt.Clusters[cID]; !ok {
		// initialise if not exists
		list = map[string]struct{}{}
	}

	payload := serviceHelper.DecodeService(bytes)

	for _, val := range payload.Data {
		if !(len(val.ID) > 0) {
			return restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
		}
		list[val.ID] = exists
	}

	cdnt.Clusters[cID] = list

	for key, _ := range list {
		services[key] = cdnt.Services[key]
	}

	// return cluster services
	payload = &service.ServicePayload{
		Data:     serviceHelper.MarshalServiceToStandardResource(services),
		Included: []service.ServiceResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}

	return serviceHelper.SendServiceWith(w, payload, http.StatusOK)
}

// Clean up the register list.
func (cdnt *Coordinator) Clean() {
	services := cdnt.Services

	for _, s := range services {
		var (
			regs       = s.RegList
			heartbeats = s.Heartbeats
		)

		for _, r := range regs {
			if t, ok := heartbeats[r.GetFullEndPoint()]; !ok ||
				time.Now().Unix()-t > int64(constants.DefaultCollaboratorExpiryInterval) {
				s.DeRegister(&r)
			}
		}
	}
}

func (cdnt *Coordinator) Dump() {
	// dump registry data to local dat file
	// constants.DataStorePath
	err1 := os.Chmod(constants.DefaultDataStorePath, 0777)
	if err1 != nil {
		logger.LogWarning("Create new data source.")
		logger.GetLoggerInstance().LogWarning("Create new data source.")
	}
	mal, err2 := json.Marshal(&cdnt)
	err2 = ioutil.WriteFile(constants.DefaultDataStorePath, mal, os.ModeExclusive|os.ModeAppend)
	if err2 != nil {
		panic(err2)
	}
	return
}
