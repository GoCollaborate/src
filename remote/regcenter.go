package remote

import (
	"encoding/json"
	"fmt"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/restful"
	"github.com/GoCollaborate/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"runtime"
	"time"
)

const (
	RPCTokenLength   = 10
	RPCServiceLength = 15
)

type RegCenter struct {
	Created  int64               // time in epoch milliseconds indicating when the registry center was created
	Modified int64               // time in epoch milliseconds indicating when the registry center was last modified
	Services map[string]*Service // map of ServiceID to Services
	Port     int
	GoVer    string
	Logger   logger.Logger
}

func NewRegCenter(port int, lg *logger.Logger) *RegCenter {
	return &RegCenter{time.Now().Unix(), time.Now().Unix(), map[string]*Service{}, port, runtime.Version(), *lg}
}

func (rc *RegCenter) Handle(router *mux.Router) *mux.Router {
	router.HandleFunc("/", handlerFuncRegCenter)
	return router
}

// post
func (rc *RegCenter) CreateService(w http.ResponseWriter, js string) (string, error) {
	// parse request json and create service here
	payload := restful.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			id := utils.RandStringBytesMaskImprSrc(RPCServiceLength)
			svrs := Service{id, dat.Attributes["description"].(string),
				dat.Attributes["parameters"].([]Parameter),
				dat.Attributes["registers"].([]Agent),
				dat.Attributes["subscribers"].([]string),
				dat.Attributes["dependencies"].([]string),
				dat.Attributes["mode"].(mode),
				dat.Attributes["load_balance_mode"].(mode),
				dat.Attributes["version"].(string),
				dat.Attributes["platform_version"].(string),
				"",
				time.Now().Unix()}
			rc.Services[id] = &svrs
			dat.ID = id
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

// post
func (rc *RegCenter) RegisterService(w http.ResponseWriter, js string) (string, error) {
	payload := restful.Decode(js)
	for _, dat := range payload.Data {
		if dat.Type == "service" {
			agent := dat.Attributes["agent"].(Agent)
			sid := dat.Attributes["serviceid"].(string)
			err := rc.Services[sid].Register(&agent)
			if err != nil {
				errPayload := restful.Error404NotFound()
				mal, er := json.Marshal(errPayload)
				if er != nil {
					return "", er
				}
				utils.AdaptHTTPWithHeader(w, constants.Header404NotFound)
				rtjs := string(mal)
				io.WriteString(w, rtjs)
				return sid, err
			}
		}
	}
	utils.AdaptHTTPWithHeader(w, constants.Header200OK)
	rtjs := string("")
	io.WriteString(w, rtjs)
	return rtjs, nil
}

// post
func (rc *RegCenter) SubscribeService(w http.ResponseWriter, srvID string) (string, error) {
	var tk string
	if _, ok := rc.Services[srvID]; ok {
		tk = utils.RandStringBytesMaskImprSrc(RPCTokenLength)
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
		utils.AdaptHTTPWithHeader(w, constants.Header200OK)
		rtjs := string("")
		io.WriteString(w, rtjs)
		return tk, nil
	}
	return tk, constants.ErrNoService
}

// delete
func (rc *RegCenter) DeleteService(w http.ResponseWriter, srvID string) (string, error) {
	return "", nil
}

// delete
func (rc *RegCenter) DeRegisterService(w http.ResponseWriter, js string) (string, error) {
	// js comprises service id and agent
	return "", nil
}

// delete
func (rc *RegCenter) UnSubscribeService(w http.ResponseWriter, js string) (string, error) {
	// js comprises service id and token
	return "", nil
}

// put
func (rc *RegCenter) AlterServiceMode(srvID string, md mode) (string, error) {
	return "", nil
}

func (rc *RegCenter) HeartBeatAll() {
	// send heartbeat signals to all service providers
	// return which of them have already expired
	return
}

func (rc *RegCenter) RefreshList() {
	rc.HeartBeatAll()
	// acquire logs from heartbeat detection and remove those which have already expired
	return
}

func (rc *RegCenter) Dump() {
	// dump registry data to local dat file
	// constants.DataStorePath
	return
}

func handlerFuncRegCenter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RegCenter Println...")
	return
}
