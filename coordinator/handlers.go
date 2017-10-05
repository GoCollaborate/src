package coordinator

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/card"
	"github.com/GoCollaborate/artifacts/restful"
	"github.com/GoCollaborate/artifacts/service"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/utils"
	"github.com/GoCollaborate/wrappers/restfulHelper"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// GET
func HandlerFuncGetServices(w http.ResponseWriter, r *http.Request) {
	// look up service list, returning service providers
	payload := restful.ExposurePayload{
		restfulHelper.MarshalServiceToStandardResource(center.Services),
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
func HandlerFuncGetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	if _, ok := center.Services[srvID]; ok {
		services := map[string]*service.Service{}
		services[srvID] = center.Services[srvID]
		payload := restful.ExposurePayload{
			restfulHelper.MarshalServiceToStandardResource(services),
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

// DELETE
func HandlerFuncDeleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	center.DeleteService(w, r, srvID)
}

// POST
func HandlerFuncPostServices(w http.ResponseWriter, r *http.Request) {
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
func HandlerFuncRegisterService(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)
	srvID := vars["srvid"]

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
	center.RegisterService(w, string(bytes), srvID)
	return
}

// POST
func HandlerFuncSubscribeService(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)
	srvID := vars["srvid"]

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
	center.SubscribeService(w, string(bytes), srvID)
	return
}

// DELETE
func HandlerFuncDeRegisterServiceAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	center.DeRegisterServiceAll(w, r, srvID)
	return
}

// DELETE
func HandlerFuncUnSubscribeServiceAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	center.UnSubscribeServiceAll(w, r, srvID)
	return
}

// DELETE
func HandlerFuncDeRegisterService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	ip := vars["ip"]
	port, _ := strconv.Atoi(vars["port"])
	center.DeRegisterService(w, r, srvID, &card.Card{ip, port, true, ""})
	return
}

// DELETE
func HandlerFuncUnSubscribeService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	token := vars["token"]
	center.UnSubscribeService(w, r, srvID, token)
	return
}

// return error if client is not in the subscriber list
func HandlerFuncQueryService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srvID := vars["srvid"]
	token := vars["token"]
	center.QueryService(w, r, srvID, token)
}
