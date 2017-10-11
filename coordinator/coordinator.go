package coordinator

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/card"
	"github.com/GoCollaborate/artifacts/restful"
	"github.com/GoCollaborate/artifacts/service"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/utils"
	"github.com/GoCollaborate/wrappers/restfulHelper"
	"github.com/GoCollaborate/wrappers/serviceHelper"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	TokenLength   = 10
	ServiceLength = 15
)

type Coordinator struct {
	Created  int64                       `json:"created"`  // time in epoch milliseconds indicating when the registry center was created
	Modified int64                       `json:"modified"` // time in epoch milliseconds indicating when the registry center was last modified
	Services map[string]*service.Service `json:"services"` // map of ServiceID to Services
	Port     int                         `json:"port"`
	GoVer    string                      `json:"gover"`
}

func (cdnt *Coordinator) Handle(router *mux.Router) *mux.Router {

	// router.HandleFunc("/", web.Index).Methods("GET")
	// router.HandleFunc("/index", web.Index).Methods("GET")
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(constants.LibUnixDir+"static/"))))
	router.HandleFunc("/services", HandlerFuncGetServices).Methods("GET")

	// client launch request to call a service
	router.HandleFunc("/query/{srvid}/{token}", HandlerFuncQueryService).Methods("GET")

	router.HandleFunc("/services/{srvid}", HandlerFuncGetService).Methods("GET")
	router.HandleFunc("/services/{srvid}", HandlerFuncDeleteService).Methods("DELETE")

	router.HandleFunc("/services/{srvid}/registry", HandlerFuncRegisterService).Methods("POST")
	router.HandleFunc("/services/{srvid}/subscription", HandlerFuncSubscribeService).Methods("POST")

	// single deletion [DELETE BODY is not explicitly specified as in HTTP 1.1,
	// so it's better not to use its message-body]
	router.HandleFunc("/services/{srvid}/registry/single/{ip}/{port}", HandlerFuncDeRegisterService).Methods("DELETE")
	router.HandleFunc("/services/{srvid}/subscription/single/{token}", HandlerFuncUnSubscribeService).Methods("DELETE")

	// bulk deletion
	router.HandleFunc("/services/{srvid}/registry", HandlerFuncDeRegisterServiceAll).Methods("DELETE")
	router.HandleFunc("/services/{srvid}/subscription", HandlerFuncUnSubscribeServiceAll).Methods("DELETE")

	// create services
	router.HandleFunc("/services", HandlerFuncPostServices).Methods("POST")
	cdnt.Dump()
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
		return serviceHelper.SendServiceWith(w, &payload, constants.Header200OK)
	}
	return restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
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
	return serviceHelper.SendServiceWith(w, &payload, constants.Header200OK)
}

// POST
func (cdnt *Coordinator) CreateService(w http.ResponseWriter, js string) error {
	// parse request json and create service here
	payload := serviceHelper.DecodeService(js)
	var resData []service.ServiceResource
	for _, dat := range payload.Data {
		if dat.Resource.Type == "service" {
			id := utils.RandStringBytesMaskImprSrc(ServiceLength)
			svrs := dat.Attributes
			svrs.ServiceID = id
			svrs.LastAssignedTime = int64(0)
			cdnt.Services[id] = &svrs
			dat.Resource.ID = id
			resData = append(resData, dat)
		}
	}
	payload.Data = resData
	serviceHelper.SendServiceWith(w, payload, constants.Header201Created)
	return nil
}

// POST
func (cdnt *Coordinator) RegisterService(w http.ResponseWriter, js string, srvID string) error {
	payload := serviceHelper.DecodeRegistry(js)
	for _, dat := range payload.Data {
		if dat.Resource.Type == "registry" {
			Cards := dat.Attributes.Cards
			if _, ok := cdnt.Services[srvID]; !ok {
				restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
				return constants.ErrNoService
			}
			for _, Card := range Cards {
				err := cdnt.Services[srvID].Register(&Card)
				if err != nil {
					restfulHelper.SendErrorWith(w, restful.Error409Conflict(), constants.Header409Conflict)
					return err
				}
			}
		}
	}
	serviceHelper.SendRegistryWith(w, payload, constants.Header200OK)
	return nil
}

// POST
func (cdnt *Coordinator) SubscribeService(w http.ResponseWriter, js string, srvID string) error {
	var res []service.SubscriptionResource
	payload := serviceHelper.DecodeSubscription(js)
	for _, dat := range payload.Data {
		if dat.Resource.Type == "subscription" {
			if _, ok := cdnt.Services[srvID]; ok {
				tk := utils.RandStringBytesMaskImprSrc(TokenLength)
				err := cdnt.Services[srvID].Subscribe(tk)
				if err != nil {
					restfulHelper.SendErrorWith(w, restful.Error409Conflict(), constants.Header409Conflict)
					return err
				}
				dat.Attributes.Token = tk
				res = append(res, dat)
			}
		}
	}
	payload.Data = res
	serviceHelper.SendSubscriptionWith(w, payload, constants.Header200OK)
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
		serviceHelper.SendServiceWith(w, &payload, constants.Header200OK)
		return nil
	}
	restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
	return constants.ErrNoService

}

// DELETE
func (cdnt *Coordinator) DeRegisterService(w http.ResponseWriter, r *http.Request, srvID string, Card *card.Card) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
		return constants.ErrNoService
	}

	err := cdnt.Services[srvID].DeRegister(Card)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
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

	serviceHelper.SendServiceWith(w, &payload, constants.Header200OK)
	return nil
}

// DELETE
func (cdnt *Coordinator) DeRegisterServiceAll(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
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

	serviceHelper.SendServiceWith(w, &payload, constants.Header200OK)
	return nil
}

// DELETE
func (cdnt *Coordinator) UnSubscribeService(w http.ResponseWriter, r *http.Request, srvID string, token string) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
		return constants.ErrNoService
	}

	err := cdnt.Services[srvID].UnSubscribe(token)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
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

	serviceHelper.SendServiceWith(w, &payload, constants.Header200OK)
	return nil
}

// DELETE
func (cdnt *Coordinator) UnSubscribeServiceAll(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
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
	serviceHelper.SendServiceWith(w, &payload, constants.Header200OK)
	return nil
}

// PUT
func (cdnt *Coordinator) AlterService(w http.ResponseWriter, js string) error {
	payload := serviceHelper.DecodeService(js)
	for _, dat := range payload.Data {
		if dat.Resource.Type == "service" {
			srvID := dat.Resource.ID
			if _, ok := cdnt.Services[srvID]; !ok {
				restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
				return constants.ErrNoService
			}
			cdnt.Services[srvID] = &dat.Attributes
		}
	}
	serviceHelper.SendServiceWith(w, payload, constants.Header202Accepted)
	return nil
}

func (cdnt *Coordinator) QueryService(w http.ResponseWriter, r *http.Request, srvID string, token string) error {
	if _, ok := cdnt.Services[srvID]; !ok {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
		return constants.ErrNoService
	}

	agt, err := cdnt.Services[srvID].Query(token)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), constants.Header404NotFound)
		return err
	}

	payload := service.QueryPayload{
		Data:     serviceHelper.MarshalCardToStandardResource(srvID, agt),
		Included: []service.QueryResource{},
		Links: restful.Links{
			Self: r.URL.String(),
		},
	}

	serviceHelper.SendQueryWith(w, &payload, constants.Header200OK)
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

func (rc *Coordinator) RefreshList() {
	// acquire logs from heartbeat detection and remove those which have already expired
	for _, service := range rc.Services {
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

func (rc *Coordinator) Dump() {
	// dump registry data to local dat file
	// constants.DataStorePath
	err1 := os.Chmod(constants.DefaultDataStorePath, 0777)
	if err1 != nil {
		logger.LogWarning("Create new data source.")
		logger.GetLoggerInstance().LogWarning("Create new data source.")
	}
	mal, err2 := json.Marshal(&rc)
	err2 = ioutil.WriteFile(constants.DefaultDataStorePath, mal, os.ModeExclusive|os.ModeAppend)
	if err2 != nil {
		panic(err2)
	}
	return
}
