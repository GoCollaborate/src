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

// must be validated before request passed in
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
				dat.Attributes["load_balance_mode"].(mode), "", time.Now().Unix()}
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

func (rc *RegCenter) RegisterService() (string, error) {
	return "", nil
}

func (rc *RegCenter) SubscribeService(srvID string) (string, error) {
	var tk string
	if _, ok := rc.Services[srvID]; ok {
		tk = utils.RandStringBytesMaskImprSrc(RPCTokenLength)
		rc.Services[srvID].SbscrbList = append(rc.Services[srvID].SbscrbList, tk)
		return tk, nil
	}
	return tk, constants.ErrNoService
}

func (rc *RegCenter) DeleteService(srvID string) (string, error) {
	return "", nil
}

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
