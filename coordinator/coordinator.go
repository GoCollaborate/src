package coordinator

import (
	"encoding/json"
	"github.com/GoCollaborate/src/artifacts/service"
	"github.com/GoCollaborate/src/constants"
	"github.com/GoCollaborate/src/logger"
	"github.com/gorilla/mux"
	"io/ioutil"
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
	Port     int32                          `json:"port"`
	GoVer    string                         `json:"go_version"`
}

func (cdnt *Coordinator) Handle(router *mux.Router) *mux.Router {

	// get services
	router.HandleFunc("/{version}/services", HandlerFuncGetServices).Methods("GET")
	// create services
	router.HandleFunc("/{version}/services", HandlerFuncPostServices).Methods("POST")
	// alter services
	router.HandleFunc("/{version}/services", HandlerFuncAlterServices).Methods("PUT")
	// alter a single service
	router.HandleFunc("/{version}/services/{srvid}", HandlerFuncAlterService).Methods("PUT")
	// get a single service
	router.HandleFunc("/{version}/services/{srvid}", HandlerFuncGetService).Methods("GET")
	// delete a single service
	router.HandleFunc("/{version}/services/{srvid}", HandlerFuncDeleteService).Methods("DELETE")
	// register a provider
	router.HandleFunc("/{version}/services/{srvid}/registry", HandlerFuncRegisterService).Methods("POST")
	// subscribe a service
	router.HandleFunc("/{version}/services/{srvid}/subscription", HandlerFuncSubscribeService).Methods("POST")

	// client launch request to call a service
	router.HandleFunc("/{version}/query/{srvid}/{token}", HandlerFuncQueryService).Methods("GET")

	// single deletion [DELETE BODY is not explicitly specified in HTTP 1.1,
	// so it's better not to use its message-body]
	router.HandleFunc("/{version}/services/{srvid}/registry/{ip}/{port}", HandlerFuncDeRegisterService).Methods("DELETE")
	router.HandleFunc("/{version}/services/{srvid}/subscription/{token}", HandlerFuncUnSubscribeService).Methods("DELETE")

	// // bulk deletion
	// router.HandleFunc("/{version}/services/{srvid}/registry", HandlerFuncDeRegisterServiceAll).Methods("DELETE")
	// router.HandleFunc("/{version}/services/{srvid}/subscription", HandlerFuncUnSubscribeServiceAll).Methods("DELETE")

	// clustering apis
	router.HandleFunc("/{version}/cluster/{id}/services", HanlderFuncGetClusterServices).Methods("GET")
	router.HandleFunc("/{version}/cluster/{id}/services", HanlderFuncPostClusterServices).Methods("POST")

	// services should be decoupled from cluster
	router.HandleFunc("/{version}/heartbeat", HandlerFuncHeartbeatServices).Methods("POST")

	go func() {
		<-time.After(constants.DefaultCollaboratorExpiryInterval)
		cdnt.Clean()
		cdnt.Dump()
	}()
	return router
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
