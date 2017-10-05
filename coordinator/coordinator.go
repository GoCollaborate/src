package coordinator

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/card"
	"github.com/GoCollaborate/artifacts/restful"
	"github.com/GoCollaborate/artifacts/service"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/utils"
	"github.com/GoCollaborate/web"
	"github.com/GoCollaborate/wrappers/cardHelper"
	"github.com/GoCollaborate/wrappers/parameterHelper"
	"github.com/GoCollaborate/wrappers/restfulHelper"
	"github.com/GoCollaborate/wrappers/serviceHelper"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	RPCTokenLength   = 10
	RPCServiceLength = 15
)

var center *Coordinator
var once sync.Once

type Coordinator struct {
	Created  int64                       `json:"created"`  // time in epoch milliseconds indicating when the registry center was created
	Modified int64                       `json:"modified"` // time in epoch milliseconds indicating when the registry center was last modified
	Services map[string]*service.Service `json:"services"` // map of ServiceID to Services
	Port     int                         `json:"port"`
	GoVer    string                      `json:"gover"`
}

// singleton instance of registry center
func GetCoordinatorInstance(port int) *Coordinator {
	once.Do(func() {
		center = &Coordinator{time.Now().Unix(), time.Now().Unix(), map[string]*service.Service{}, port, runtime.Version()}
	})
	return center
}

func (rc *Coordinator) Handle(router *mux.Router) *mux.Router {

	center = rc

	router.HandleFunc("/", web.Index).Methods("GET")
	router.HandleFunc("/index", web.Index).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(constants.LibUnixDir+"static/"))))
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
	rc.Dump()
	return router
}

// POST
func (rc *Coordinator) CreateService(w http.ResponseWriter, js string) error {
	// parse request json and create service here
	payload := restfulHelper.Decode(js)
	var resData []restful.Resource
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			id := utils.RandStringBytesMaskImprSrc(RPCServiceLength)
			svrs := service.Service{id, dat.Attributes["description"].(string),
				parameterHelper.UnmarshalParameters(dat.Attributes["parameters"].([]interface{})),
				cardHelper.UnmarshalCards(dat.Attributes["registers"].([]interface{})),
				parameterHelper.UnmarshalStringArray(dat.Attributes["subscribers"].([]interface{})),
				parameterHelper.UnmarshalStringArray(dat.Attributes["dependencies"].([]interface{})),
				serviceHelper.UnmarshalMode(dat.Attributes["mode"]),
				serviceHelper.UnmarshalMode(dat.Attributes["load_balance_mode"]),
				dat.Attributes["version"].(string),
				dat.Attributes["platform_version"].(string),
				"",
				time.Now().Unix()}
			rc.Services[id] = &svrs
			dat.ID = id
			resData = append(resData, dat)
		}
	}
	payload.Data = resData
	mal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header201Created)
	io.WriteString(w, string(mal))
	return nil
}

// POST
func (rc *Coordinator) RegisterService(w http.ResponseWriter, js string, srvID string) error {
	payload := restfulHelper.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "registry" {
			Cards := cardHelper.UnmarshalCards(dat.Attributes["cards"].([]interface{}))
			if _, ok := rc.Services[srvID]; !ok {
				errPayload := restful.Error404NotFound()
				mal, er := json.Marshal(errPayload)
				if er != nil {
					return er
				}
				utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
				io.WriteString(w, string(mal))
				return constants.ErrNoService
			}
			for _, Card := range Cards {
				err := rc.Services[srvID].Register(&Card)
				if err != nil {
					errPayload := restful.Error409Conflict()
					mal, er := json.Marshal(errPayload)
					if er != nil {
						return er
					}
					utils.AdaptHTTPWithHeader(w, constants.Header409Conflict)
					rtjs := string(mal)
					io.WriteString(w, rtjs)
					return err
				}
			}
		}
	}
	mal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	io.WriteString(w, string(mal))
	return nil
}

// POST
func (rc *Coordinator) SubscribeService(w http.ResponseWriter, js string, srvID string) error {
	var tks []restful.Resource
	payload := restfulHelper.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "subscription" {
			if _, ok := rc.Services[srvID]; ok {
				tk := utils.RandStringBytesMaskImprSrc(RPCTokenLength)
				err := rc.Services[srvID].Subscribe(tk)
				if err != nil {
					errPayload := restful.Error409Conflict()
					mal, er := json.Marshal(errPayload)
					if er != nil {
						return er
					}
					utils.AdaptHTTPWithHeader(w, constants.Header409Conflict)
					rtjs := string(mal)
					io.WriteString(w, rtjs)
					return err
				}
				if dat.Attributes == nil {
					dat.Attributes = map[string]interface{}{}
				}
				dat.Attributes["token"] = tk
				tks = append(tks, dat)
			}
		}
	}
	payload.Data = tks
	mal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	io.WriteString(w, string(mal))
	return nil
}

// DELETE
func (rc *Coordinator) DeleteService(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := rc.Services[srvID]; ok {
		services := map[string]*service.Service{}
		services[srvID] = center.Services[srvID]
		payload := restful.ExposurePayload{
			restfulHelper.MarshalServiceToStandardResource(services),
			[]restful.Resource{},
			restful.Links{
				Self: r.URL.String(),
			},
		}
		delete(rc.Services, srvID)
		mal, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		utils.AdaptHTTPWithHeader(w, constants.Header200OK)
		io.WriteString(w, string(mal))
		return nil
	}
	utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
	io.WriteString(w, "")
	return constants.ErrNoService

}

// DELETE
func (rc *Coordinator) DeRegisterService(w http.ResponseWriter, r *http.Request, srvID string, Card *card.Card) error {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		io.WriteString(w, string(mal))
		return constants.ErrNoService
	}

	err := rc.Services[srvID].DeRegister(Card)
	if err != nil {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		io.WriteString(w, string(mal))
		return err
	}

	services := map[string]*service.Service{}
	services[srvID] = rc.Services[srvID]

	rtpayload := restful.ExposurePayload{
		restfulHelper.MarshalServiceToStandardResource(services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}

	mal, errr := json.Marshal(rtpayload)
	if errr != nil {
		return errr
	}

	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	io.WriteString(w, string(mal))
	return nil
}

// DELETE
func (rc *Coordinator) DeRegisterServiceAll(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		io.WriteString(w, string(mal))
		return constants.ErrNoService
	}
	services := map[string]*service.Service{}
	services[srvID] = rc.Services[srvID]
	payload := restful.ExposurePayload{
		restfulHelper.MarshalServiceToStandardResource(services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}
	rc.Services[srvID].DeRegisterAll()

	mal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	io.WriteString(w, string(mal))
	return nil
}

// DELETE
func (rc *Coordinator) UnSubscribeService(w http.ResponseWriter, r *http.Request, srvID string, token string) error {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, _ := json.Marshal(errPayload)

		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		io.WriteString(w, string(mal))
		return constants.ErrNoService
	}

	err := rc.Services[srvID].UnSubscribe(token)
	if err != nil {
		errPayload := restful.Error404NotFound()
		mal, _ := json.Marshal(errPayload)

		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return err
	}

	services := map[string]*service.Service{}
	services[srvID] = rc.Services[srvID]

	rtpayload := restful.ExposurePayload{
		restfulHelper.MarshalServiceToStandardResource(services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}

	mal, errr := json.Marshal(rtpayload)
	if errr != nil {
		return errr
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	io.WriteString(w, string(mal))
	return nil
}

// DELETE
func (rc *Coordinator) UnSubscribeServiceAll(w http.ResponseWriter, r *http.Request, srvID string) error {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		io.WriteString(w, string(mal))
		return constants.ErrNoService
	}
	services := map[string]*service.Service{}
	services[srvID] = rc.Services[srvID]
	payload := restful.ExposurePayload{
		restfulHelper.MarshalServiceToStandardResource(services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}
	rc.Services[srvID].UnSubscribeAll()
	mal, _ := json.Marshal(payload)

	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	io.WriteString(w, string(mal))
	return nil
}

// PUT
func (rc *Coordinator) AlterService(w http.ResponseWriter, js string) error {
	payload := restfulHelper.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			srvID := dat.ID
			if _, ok := rc.Services[srvID]; !ok {
				errPayload := restful.Error404NotFound()
				mal, er := json.Marshal(errPayload)
				if er != nil {
					return er
				}
				utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
				io.WriteString(w, string(mal))
				return constants.ErrNoService
			}

			svrs := service.Service{dat.ID, dat.Attributes["description"].(string),
				parameterHelper.UnmarshalParameters(dat.Attributes["parameters"].([]interface{})),
				cardHelper.UnmarshalCards(dat.Attributes["registers"].([]interface{})),
				parameterHelper.UnmarshalStringArray(dat.Attributes["subscribers"].([]interface{})),
				parameterHelper.UnmarshalStringArray(dat.Attributes["dependencies"].([]interface{})),
				serviceHelper.UnmarshalMode(dat.Attributes["mode"]),
				serviceHelper.UnmarshalMode(dat.Attributes["load_balance_mode"]),
				dat.Attributes["version"].(string),
				dat.Attributes["platform_version"].(string),
				"",
				time.Now().Unix()}
			rc.Services[srvID] = &svrs
			dat.ID = srvID
		}
	}
	mal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header202Accepted)
	io.WriteString(w, string(mal))
	return nil
}

func (rc *Coordinator) QueryService(w http.ResponseWriter, r *http.Request, srvID string, token string) error {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		io.WriteString(w, string(mal))
		return constants.ErrNoService
	}

	agt, err := rc.Services[srvID].Query(token)
	if err != nil {
		errPayload := restful.Error404NotFound()
		mal, _ := json.Marshal(errPayload)

		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return err
	}

	rtpayload := restful.ExposurePayload{
		restfulHelper.MarshalCardToStandardResource(srvID, agt),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}

	mal, _ := json.Marshal(rtpayload)

	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	io.WriteString(w, string(mal))
	return nil
}

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
