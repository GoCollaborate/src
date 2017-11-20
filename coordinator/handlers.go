package coordinator

import (
	"github.com/GoCollaborate/src/artifacts/card"
	"github.com/GoCollaborate/src/artifacts/restful"
	"github.com/GoCollaborate/src/artifacts/service"
	"github.com/GoCollaborate/src/wrappers/restfulHelper"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var singleton *Coordinator
var once sync.Once

// singleton instance of registry center
func GetCoordinatorInstance(port int32) *Coordinator {
	once.Do(func() {
		singleton = &Coordinator{time.Now().Unix(), time.Now().Unix(), map[string]*service.Service{}, map[string]map[string]struct{}{}, port, runtime.Version()}
	})
	return singleton
}

// GET
func HandlerFuncGetServices(w http.ResponseWriter, r *http.Request) {
	singleton.GetServices(w, r)
	return
}

// GET
func HandlerFuncGetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	singleton.GetService(w, r, srvID)
	return
}

// DELETE
func HandlerFuncDeleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	singleton.DeleteService(w, r, srvID)
}

// POST
func HandlerFuncPostServices(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error415UnsupportedMediaType(), http.StatusUnsupportedMediaType)
		return
	}
	singleton.CreateService(w, bytes)
	return
}

// PUT
func HandlerFuncAlterServices(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error415UnsupportedMediaType(), http.StatusUnsupportedMediaType)
		return
	}
	singleton.AlterService(w, bytes)
	return
}

// POST
func HandlerFuncRegisterService(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)
	srvID := vars["srvid"]

	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error415UnsupportedMediaType(), http.StatusUnsupportedMediaType)
		return
	}
	singleton.RegisterService(w, bytes, srvID)
	return
}

// POST
func HandlerFuncSubscribeService(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)
	srvID := vars["srvid"]

	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error415UnsupportedMediaType(), http.StatusUnsupportedMediaType)
		return
	}
	singleton.SubscribeService(w, bytes, srvID)
	return
}

// DELETE
func HandlerFuncDeRegisterServiceAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	singleton.DeRegisterServiceAll(w, r, srvID)
	return
}

// DELETE
func HandlerFuncUnSubscribeServiceAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	singleton.UnSubscribeServiceAll(w, r, srvID)
	return
}

// DELETE
func HandlerFuncDeRegisterService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	ip := vars["ip"]
	port, _ := strconv.Atoi(vars["port"])
	singleton.DeRegisterService(w, r, srvID, &card.Card{ip, int32(port), true, "", false})
	return
}

// DELETE
func HandlerFuncUnSubscribeService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	token := vars["token"]
	singleton.UnSubscribeService(w, r, srvID, token)
	return
}

// return error if client is not in the subscriber list
func HandlerFuncQueryService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	token := vars["token"]
	singleton.QueryService(w, r, srvID, token)
}

// a handler to maintain service provider's heartbeats,
// this interface will most likely be changed to protobuf message in the future
func HandlerFuncHeartbeatServices(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error415UnsupportedMediaType(), http.StatusUnsupportedMediaType)
		return
	}
	singleton.HeartbeatServices(w, bytes)
}

func HanlderFuncGetClusterServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cID := vars["cid"]
	singleton.GetClusterServices(w, r, cID)
}

func HanlderFuncPostClusterServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cID := vars["cid"]

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		restfulHelper.SendErrorWith(w, restful.Error415UnsupportedMediaType(), http.StatusUnsupportedMediaType)
		return
	}

	singleton.PostClusterServices(w, r, cID, bytes)
}
