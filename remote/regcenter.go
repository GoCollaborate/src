package remote

import (
	"encoding/base64"
	"encoding/json"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/restful"
	"github.com/GoCollaborate/utils"
	"github.com/GoCollaborate/web"
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

var center *RegCenter
var once sync.Once

type RegCenter struct {
	Created  int64               `json:"created"`  // time in epoch milliseconds indicating when the registry center was created
	Modified int64               `json:"modified"` // time in epoch milliseconds indicating when the registry center was last modified
	Services map[string]*Service `json:"services"` // map of ServiceID to Services
	Port     int                 `json:"port"`
	GoVer    string              `json:"gover"`
	Logger   logger.Logger       `json:"-"`
}

// singleton instance of registry center
func GetRegCenterInstance(port int, lg *logger.Logger) *RegCenter {
	once.Do(func() {
		center = &RegCenter{time.Now().Unix(), time.Now().Unix(), map[string]*Service{}, port, runtime.Version(), *lg}
	})
	return center
}

func (rc *RegCenter) Handle(router *mux.Router) *mux.Router {

	center = rc

	router.HandleFunc("/", web.Index).Methods("GET")
	router.HandleFunc("/index", web.Index).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	router.HandleFunc("/services", handlerFuncGetServices).Methods("GET")
	router.HandleFunc("/services/{srvid}", handlerFuncGetService).Methods("GET")
	// todo: service specific api entry e.g. /services/{srvid}/registry, no id field required
	router.HandleFunc("/services/registry", handlerFuncRegisterService).Methods("POST")
	router.HandleFunc("/services/subscription", handlerFuncSubscribeService).Methods("POST")

	// with http parameters, if the parameters contains url, it should be encoded in base64 format
	router.HandleFunc("/services/{srvid}/registry/{parjson}", handlerFuncDeRegisterService).Methods("DELETE")
	router.HandleFunc("/services/{srvid}/subscription/{token}", handlerFuncUnSubscribeService).Methods("DELETE")

	// without http parameters
	router.HandleFunc("/services/{srvid}/registry", handlerFuncDeRegisterServiceAll).Methods("DELETE")
	router.HandleFunc("/services/{srvid}/subscription", handlerFuncUnSubscribeServiceAll).Methods("DELETE")
	router.HandleFunc("/services", handlerFuncPostServices).Methods("POST")
	rc.Dump()
	return router
}

// POST
func (rc *RegCenter) CreateService(w http.ResponseWriter, js string) (string, error) {
	// parse request json and create service here
	payload := restful.Decode(js)
	var resData []restful.Resource
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			id := utils.RandStringBytesMaskImprSrc(RPCServiceLength)
			svrs := Service{id, dat.Attributes["description"].(string),
				UnmarshalParameters(dat.Attributes["parameters"].([]interface{})),
				UnmarshalAgents(dat.Attributes["registers"].([]interface{})),
				UnmarshalStringArray(dat.Attributes["subscribers"].([]interface{})),
				UnmarshalStringArray(dat.Attributes["dependencies"].([]interface{})),
				UnmarshalMode(dat.Attributes["mode"]),
				UnmarshalMode(dat.Attributes["load_balance_mode"]),
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
		return "", err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header201Created)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return rtjs, nil
}

// POST
func (rc *RegCenter) RegisterService(w http.ResponseWriter, js string) (string, error) {
	payload := restful.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "registry" {
			agents := UnmarshalAgents(dat.Attributes["agents"].([]interface{}))
			srvID := dat.ID
			if _, ok := rc.Services[srvID]; !ok {
				errPayload := restful.Error404NotFound()
				mal, er := json.Marshal(errPayload)
				if er != nil {
					return "", er
				}
				utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
				rtjs := string(mal)
				io.WriteString(w, rtjs)
				return srvID, constants.ErrNoService
			}
			for _, agent := range agents {
				err := rc.Services[srvID].Register(&agent)
				if err != nil {
					errPayload := restful.Error409Conflict()
					mal, er := json.Marshal(errPayload)
					if er != nil {
						return "", er
					}
					utils.AdaptHTTPWithHeader(w, constants.Header409Conflict)
					rtjs := string(mal)
					io.WriteString(w, rtjs)
					return srvID, err
				}
			}
		}
	}
	mal, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return rtjs, nil
}

// POST
func (rc *RegCenter) SubscribeService(w http.ResponseWriter, js string) (string, error) {
	var tks []restful.Resource
	payload := restful.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "subscription" {
			srvID := dat.ID
			if _, ok := rc.Services[srvID]; ok {
				tk := utils.RandStringBytesMaskImprSrc(RPCTokenLength)
				err := rc.Services[srvID].Subscribe(tk)
				if err != nil {
					errPayload := restful.Error409Conflict()
					mal, er := json.Marshal(errPayload)
					if er != nil {
						return "", er
					}
					utils.AdaptHTTPWithHeader(w, constants.Header409Conflict)
					rtjs := string(mal)
					io.WriteString(w, rtjs)
					return tk, err
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
		return "", err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return "", nil
}

// DELETE
func (rc *RegCenter) DeleteService(w http.ResponseWriter, js string) (string, error) {
	payload := restful.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			srvID := dat.ID
			if _, ok := rc.Services[srvID]; ok {
				delete(rc.Services, srvID)
				continue
			}
			utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
			rtjs := string("")
			io.WriteString(w, rtjs)
			return srvID, constants.ErrNoService
		}
	}
	mal, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return "", nil
}

// DELETE
func (rc *RegCenter) DeRegisterService(w http.ResponseWriter, r *http.Request, srvID string, agent *Agent) (string, error) {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return "", er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return srvID, constants.ErrNoService
	}

	services := map[string]*Service{}
	services[srvID] = rc.Services[srvID]

	err := rc.Services[srvID].DeRegister(agent)
	if err != nil {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return "", er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return srvID, err
	}
	payload := restful.Payload{
		marshalServiceToStandardResource(services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}

	mal, errr := json.Marshal(payload)
	if errr != nil {
		return "", errr
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return srvID, nil
}

// DELETE
func (rc *RegCenter) DeRegisterServiceAll(w http.ResponseWriter, r *http.Request, srvID string) (string, error) {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return "", er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return srvID, constants.ErrNoService
	}
	services := map[string]*Service{}
	services[srvID] = rc.Services[srvID]
	payload := restful.Payload{
		marshalServiceToStandardResource(services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}
	rc.Services[srvID].DeRegisterAll()

	mal, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return rtjs, nil
}

// DELETE
func (rc *RegCenter) UnSubscribeService(w http.ResponseWriter, r *http.Request, srvID string, token string) (string, error) {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return "", er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return srvID, constants.ErrNoService
	}

	services := map[string]*Service{}
	services[srvID] = rc.Services[srvID]

	err := rc.Services[srvID].UnSubscribe(token)
	if err != nil {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return "", er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return srvID, err
	}
	payload := restful.Payload{
		marshalServiceToStandardResource(services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}

	mal, errr := json.Marshal(payload)
	if errr != nil {
		return "", errr
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return srvID, nil
}

// DELETE
func (rc *RegCenter) UnSubscribeServiceAll(w http.ResponseWriter, r *http.Request, srvID string) (string, error) {
	if _, ok := rc.Services[srvID]; !ok {
		errPayload := restful.Error404NotFound()
		mal, er := json.Marshal(errPayload)
		if er != nil {
			return "", er
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return srvID, constants.ErrNoService
	}
	services := map[string]*Service{}
	services[srvID] = rc.Services[srvID]
	payload := restful.Payload{
		marshalServiceToStandardResource(services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}
	rc.Services[srvID].UnSubscribeAll()
	mal, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return "", nil
}

// PUT
func (rc *RegCenter) AlterService(w http.ResponseWriter, js string) (string, error) {
	payload := restful.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			srvID := dat.ID
			if _, ok := rc.Services[srvID]; !ok {
				errPayload := restful.Error404NotFound()
				mal, er := json.Marshal(errPayload)
				if er != nil {
					return "", er
				}
				utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
				rtjs := string(mal)
				io.WriteString(w, rtjs)
				return srvID, constants.ErrNoService
			}

			svrs := Service{dat.ID, dat.Attributes["description"].(string),
				UnmarshalParameters(dat.Attributes["parameters"].([]interface{})),
				UnmarshalAgents(dat.Attributes["registers"].([]interface{})),
				UnmarshalStringArray(dat.Attributes["subscribers"].([]interface{})),
				UnmarshalStringArray(dat.Attributes["dependencies"].([]interface{})),
				UnmarshalMode(dat.Attributes["mode"]),
				UnmarshalMode(dat.Attributes["load_balance_mode"]),
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
		return "", err
	}
	utils.AdaptHTTPWithHeader(w, constants.Header202Accepted)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return rtjs, nil
}

func (rc *RegCenter) HeartBeatAll() []Agent {
	// send heartbeat signals to all service providers
	// return which of them have already expired
	var expired []Agent
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

func (rc *RegCenter) RefreshList() {
	// acquire logs from heartbeat detection and remove those which have already expired
	for _, service := range rc.Services {
		var activated []Agent
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

func (rc *RegCenter) Dump() {
	// dump registry data to local dat file
	// constants.DataStorePath
	err1 := os.Chmod(constants.DefaultDataStorePath, 0777)
	if err1 != nil {
		panic(err1)
	}
	mal, err2 := json.Marshal(&rc)
	err2 = ioutil.WriteFile(constants.DefaultDataStorePath, mal, os.ModeExclusive)
	if err2 != nil {
		panic(err2)
	}
	return
}

// GET
func handlerFuncGetServices(w http.ResponseWriter, r *http.Request) {
	// look up service list, returning service providers
	payload := restful.Payload{
		marshalServiceToStandardResource(center.Services),
		[]restful.Resource{},
		restful.Links{
			Self: r.URL.String(),
		},
	}
	mal, err := json.Marshal(payload)
	if err != nil {
		return
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return
}

// GET
func handlerFuncGetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	if _, ok := center.Services[srvID]; ok {
		services := map[string]*Service{}
		services[srvID] = center.Services[srvID]
		payload := restful.Payload{
			marshalServiceToStandardResource(services),
			[]restful.Resource{},
			restful.Links{
				Self: r.URL.String(),
			},
		}
		mal, err := json.Marshal(payload)
		if err != nil {
			return
		}
		utils.AdaptHTTPWithHeader(w, constants.Header200OK)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return
	}
	errPayload := restful.Error404NotFound()
	mal, err := json.Marshal(errPayload)
	if err != nil {
		return
	}
	utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
	rtjs := string(mal)
	io.WriteString(w, rtjs)
	return
}

// POST
func handlerFuncPostServices(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errPayload := restful.Error404NotFound()
		mal, err := json.Marshal(errPayload)
		if err != nil {
			return
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return
	}
	center.CreateService(w, string(bytes))
	return
}

// POST
func handlerFuncRegisterService(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errPayload := restful.Error404NotFound()
		mal, err := json.Marshal(errPayload)
		if err != nil {
			return
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return
	}
	center.RegisterService(w, string(bytes))
	return
}

// POST
func handlerFuncSubscribeService(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errPayload := restful.Error404NotFound()
		mal, err := json.Marshal(errPayload)
		if err != nil {
			return
		}
		utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
		rtjs := string(mal)
		io.WriteString(w, rtjs)
		return
	}
	center.SubscribeService(w, string(bytes))
	return
}

// DELETE
func handlerFuncDeRegisterServiceAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	center.DeRegisterServiceAll(w, r, srvID)
	return
}

// DELETE
func handlerFuncUnSubscribeServiceAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	center.UnSubscribeServiceAll(w, r, srvID)
	return
}

// DELETE
func handlerFuncDeRegisterService(w http.ResponseWriter, r *http.Request) {
	// srvid=xxx&ip=xxx&port=xxx&alive=xxx
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	parJSONStr := vars["parjson"]
	agent, err := unMarshalRegParJSON(parJSONStr)
	if err != nil {
		return
	}
	center.DeRegisterService(w, r, srvID, &agent)
	return
}

// DELETE
func handlerFuncUnSubscribeService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	token := vars["token"]
	center.UnSubscribeService(w, r, srvID, token)
	return
}

func marshalServiceToStandardResource(srvs map[string]*Service) []restful.Resource {
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

func unMarshalRegParJSON(bs64 string) (Agent, error) {
	decoded, err := base64.StdEncoding.DecodeString(bs64)
	if err != nil {
		return Agent{}, err
	}
	var agt Agent
	err = json.Unmarshal([]byte(decoded), &agt)
	if err != nil {
		return Agent{}, err
	}
	return agt, nil
}
