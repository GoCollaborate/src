package coordinator

import (
	"github.com/GoCollaborate/src/artifacts/card"
	"github.com/GoCollaborate/src/artifacts/resources"
	"github.com/GoCollaborate/src/artifacts/restful"
	. "github.com/GoCollaborate/src/artifacts/service"
	"github.com/GoCollaborate/src/constants"
	"github.com/GoCollaborate/src/utils"
	"github.com/gorilla/mux"
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
		singleton = &Coordinator{
			time.Now().Unix(),
			time.Now().Unix(),
			map[string]*Service{},
			map[string]map[string]struct{}{},
			port,
			runtime.Version(),
		}
	})
	return singleton
}

// GET
func HandlerFuncGetService(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		srvID = vars["srvid"]
	)
	// look up service list, returning service providers
	if _, ok := singleton.Services[srvID]; ok {
		services := map[string]*Service{}
		services[srvID] = singleton.Services[srvID]

		restful.
			Writer(w).
			WithResources(
				NewServicesResource(&services),
			).
			WithLinks(
				&resources.Links{
					Self: r.URL.String(),
				},
			).
			WithHeader(
				constants.HeaderContentTypeJSON.Key,
				constants.HeaderContentTypeJSON.Value,
			).
			WithStatus(http.StatusOK).
			Send()
		return
	}
	restful.
		Writer(w).
		WithResources(
			restful.NewErrorResource(restful.Error404NotFound()),
		).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusNotFound).
		Send()
	return
}

// GET
func HandlerFuncGetServices(w http.ResponseWriter, r *http.Request) {
	restful.
		Writer(w).
		WithResources(
			NewServicesResource(
				&singleton.Services,
			),
		).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()
	return
}

// POST
func HandlerFuncPostServices(w http.ResponseWriter, r *http.Request) {
	var (
		serviceResources = new([]ServiceResource)
	)

	restful.
		Reader(r).
		WithResourceArr(serviceResources).
		Receive()

	for _, dat := range *serviceResources {
		if dat.Type == "service" {
			id := utils.RandStringBytesMaskImprSrc(ServiceLength)
			svr := NewServiceFrom(dat.Attributes)
			svr.ServiceID = id
			singleton.Services[id] = svr
			dat.Id = id
			dat.Attributes = svr
		}
	}

	restful.
		Writer(w).
		WithResourceArr(serviceResources).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()

	return
}

// DELETE
func HandlerFuncDeleteService(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		srvID = vars["srvid"]
	)

	if _, ok := singleton.Services[srvID]; ok {
		services := map[string]*Service{}
		services[srvID] = singleton.Services[srvID]
		delete(singleton.Services, srvID)
		restful.
			Writer(w).
			WithResources(
				NewServicesResource(&services),
			).
			WithLinks(
				&resources.Links{
					Self: r.URL.String(),
				},
			).
			WithHeader(
				constants.HeaderContentTypeJSON.Key,
				constants.HeaderContentTypeJSON.Value,
			).
			WithStatus(http.StatusOK).
			Send()
		return
	}

	restful.
		Writer(w).
		WithResources(
			restful.NewErrorResource(restful.Error404NotFound()),
		).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusNotFound).
		Send()
	return
}

// PUT
func HandlerFuncAlterServices(w http.ResponseWriter, r *http.Request) {

	var (
		serviceResources = new([]ServiceResource)
	)

	restful.
		Reader(r).
		WithResourceArr(serviceResources).
		Receive()

	for _, dat := range *serviceResources {
		if _, ok := singleton.Services[dat.Id]; !ok || dat.Type != "service" {

			restful.
				Writer(w).
				WithResources(
					restful.NewErrorResource(restful.Error404NotFound()),
				).
				WithLinks(
					&resources.Links{
						Self: r.URL.String(),
					},
				).
				WithHeader(
					constants.HeaderContentTypeJSON.Key,
					constants.HeaderContentTypeJSON.Value,
				).
				WithStatus(http.StatusNotFound).
				Send()
			return
		}
		singleton.Services[dat.Id] = dat.Attributes
	}

	restful.
		Writer(w).
		WithResourceArr(serviceResources).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()

	return
}

// PUT
func HandlerFuncAlterService(w http.ResponseWriter, r *http.Request) {

	var (
		serviceResource = new(ServiceResource)
		vars            = mux.Vars(r)
		srvID           = vars["srvid"]
	)

	restful.
		Reader(r).
		WithResource(serviceResource).
		Receive()

	if _, ok := singleton.Services[srvID]; !ok || serviceResource.Type != "service" {

		restful.
			Writer(w).
			WithResources(
				restful.NewErrorResource(restful.Error404NotFound()),
			).
			WithLinks(
				&resources.Links{
					Self: r.URL.String(),
				},
			).
			WithHeader(
				constants.HeaderContentTypeJSON.Key,
				constants.HeaderContentTypeJSON.Value,
			).
			WithStatus(http.StatusNotFound).
			Send()
		return
	}

	singleton.Services[srvID] = serviceResource.Attributes

	restful.
		Writer(w).
		WithResource(serviceResource).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()

	return
}

// POST
func HandlerFuncRegisterService(w http.ResponseWriter, r *http.Request) {

	var (
		registryResources = new([]CardResource)
		vars              = mux.Vars(r)
		srvID             = vars["srvid"]
	)

	restful.
		Reader(r).
		WithResourceArr(registryResources).
		Receive()

	for _, dat := range *registryResources {
		if _, ok := singleton.Services[srvID]; !ok || dat.Type != "registry" {
			restful.
				Writer(w).
				WithResources(
					restful.NewErrorResource(restful.Error404NotFound()),
				).
				WithLinks(
					&resources.Links{
						Self: r.URL.String(),
					},
				).
				WithHeader(
					constants.HeaderContentTypeJSON.Key,
					constants.HeaderContentTypeJSON.Value,
				).
				WithStatus(http.StatusNotFound).
				Send()
			return
		}

		if err := singleton.Services[srvID].Register(dat.Attributes); err != nil {
			restful.
				Writer(w).
				WithResources(
					restful.NewErrorResource(restful.Error409Conflict()),
				).
				WithLinks(
					&resources.Links{
						Self: r.URL.String(),
					},
				).
				WithHeader(
					constants.HeaderContentTypeJSON.Key,
					constants.HeaderContentTypeJSON.Value,
				).
				WithStatus(http.StatusConflict).
				Send()
			return
		}
	}

	restful.
		Writer(w).
		WithResourceArr(registryResources).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()

	return
}

// POST
func HandlerFuncSubscribeService(w http.ResponseWriter, r *http.Request) {
	var (
		vars                  = mux.Vars(r)
		srvID                 = vars["srvid"]
		subscriptionResources = new([]SubscriptionResource)
	)

	restful.
		Reader(r).
		WithResourceArr(subscriptionResources).
		Receive()

	for _, dat := range *subscriptionResources {
		if _, ok := singleton.Services[srvID]; ok && dat.Type == "subscription" {

			tk := utils.RandStringBytesMaskImprSrc(TokenLength)
			if err := singleton.Services[srvID].Subscribe(tk); err != nil {
				restful.
					Writer(w).
					WithResources(
						restful.NewErrorResource(restful.Error409Conflict()),
					).
					WithLinks(
						&resources.Links{
							Self: r.URL.String(),
						},
					).
					WithHeader(
						constants.HeaderContentTypeJSON.Key,
						constants.HeaderContentTypeJSON.Value,
					).
					WithStatus(http.StatusConflict).
					Send()
				return
			}
			dat.Id = tk
		}
	}

	restful.
		Writer(w).
		WithResourceArr(subscriptionResources).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()
	return
}

// DELETE
func HandlerFuncDeRegisterService(w http.ResponseWriter, r *http.Request) {
	var (
		vars    = mux.Vars(r)
		srvID   = vars["srvid"]
		ip      = vars["ip"]
		port, _ = strconv.ParseInt(vars["port"], 10, 32)
	)

	if _, ok := singleton.Services[srvID]; !ok || singleton.Services[srvID].DeRegister(
		card.NewCard(ip, int32(port), false, "", false),
	) != nil {

		restful.
			Writer(w).
			WithResources(
				restful.NewErrorResource(restful.Error409Conflict()),
			).
			WithLinks(
				&resources.Links{
					Self: r.URL.String(),
				},
			).
			WithHeader(
				constants.HeaderContentTypeJSON.Key,
				constants.HeaderContentTypeJSON.Value,
			).
			WithStatus(http.StatusConflict).
			Send()
		return
	}

	restful.
		Writer(w).
		WithResource(
			NewServiceResource(singleton.Services[srvID]),
		).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()
	return
}

// DELETE
func HandlerFuncUnSubscribeService(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		srvID = vars["srvid"]
		token = vars["token"]
	)

	if _, ok := singleton.Services[srvID]; !ok || singleton.Services[srvID].UnSubscribe(token) != nil {

		restful.
			Writer(w).
			WithResources(
				restful.NewErrorResource(restful.Error409Conflict()),
			).
			WithLinks(
				&resources.Links{
					Self: r.URL.String(),
				},
			).
			WithHeader(
				constants.HeaderContentTypeJSON.Key,
				constants.HeaderContentTypeJSON.Value,
			).
			WithStatus(http.StatusConflict).
			Send()
		return
	}

	restful.
		Writer(w).
		WithResource(
			NewServiceResource(singleton.Services[srvID]),
		).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()
	return
}

// // DELETE
// func HandlerFuncDeRegisterServiceAll(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	srvID := vars["srvid"]

// 	serviceHelper.SendServiceWith(w, &payload, http.StatusOK)

// 	if _, ok := singleton.Services[srvID]; !ok {
// 		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
// 		return constants.ErrNoService
// 	}
// 	services := map[string]*service.Service{}
// 	services[srvID] = singleton.Services[srvID]
// 	payload := service.ServicePayload{
// 		Data:     serviceHelper.MarshalServiceToStandardResource(services),
// 		Included: []service.ServiceResource{},
// 		Links: restful.Links{
// 			Self: r.URL.String(),
// 		},
// 	}
// 	singleton.Services[srvID].DeRegisterAll()

// 	serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
// 	return
// }

// // DELETE
// func HandlerFuncUnSubscribeServiceAll(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	srvID := vars["srvid"]

// 	if _, ok := singleton.Services[srvID]; !ok {
// 		restfulHelper.SendErrorWith(w, restful.Error404NotFound(), http.StatusNotFound)
// 		return constants.ErrNoService
// 	}
// 	services := map[string]*service.Service{}
// 	services[srvID] = singleton.Services[srvID]
// 	payload := service.ServicePayload{
// 		Data:     serviceHelper.MarshalServiceToStandardResource(services),
// 		Included: []service.ServiceResource{},
// 		Links: restful.Links{
// 			Self: r.URL.String(),
// 		},
// 	}
// 	singleton.Services[srvID].UnSubscribeAll()
// 	serviceHelper.SendServiceWith(w, &payload, http.StatusOK)
// 	return
// }

// return error if client is not in the subscriber list
func HandlerFuncQueryService(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		srvID = vars["srvid"]
		token = vars["token"]
	)

	if _, ok := singleton.Services[srvID]; !ok {
		restful.Writer(w).
			WithResources(
				restful.NewErrorResource(
					restful.Error404NotFound(),
				),
			).
			WithLinks(
				&resources.Links{
					Self: r.URL.String(),
				},
			).
			WithHeader(
				constants.HeaderContentTypeJSON.Key,
				constants.HeaderContentTypeJSON.Value,
			).
			WithStatus(http.StatusNotFound).
			Send()
		return
	}

	agt, err := singleton.Services[srvID].Query(token)
	if err != nil {
		restful.Writer(w).
			WithResources(
				restful.NewErrorResource(
					restful.Error404NotFound(),
				),
			).
			WithLinks(
				&resources.Links{
					Self: r.URL.String(),
				},
			).
			WithHeader(
				constants.HeaderContentTypeJSON.Key,
				constants.HeaderContentTypeJSON.Value,
			).
			WithStatus(http.StatusNotFound).
			Send()
		return
	}

	restful.Writer(w).
		WithResource(
			NewQueryResource(agt),
		).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()
}

// a handler to maintain service provider's heartbeats,
// this interface will most likely be changed to protobuf message in the future
func HandlerFuncHeartbeatServices(w http.ResponseWriter, r *http.Request) {

	var (
		heartbeatResources = new([]CardResource)
	)

	restful.Reader(r).
		WithResourceArr(heartbeatResources).
		Receive()

	for _, dat := range *heartbeatResources {

		if _, ok := singleton.Services[dat.Id]; !ok || dat.Type != "heartbeat" {
			restful.Writer(w).
				WithResources(
					restful.NewErrorResource(
						restful.Error404NotFound(),
					),
				).
				WithLinks(
					&resources.Links{
						Self: r.URL.String(),
					},
				).
				WithHeader(
					constants.HeaderContentTypeJSON.Key,
					constants.HeaderContentTypeJSON.Value,
				).
				WithStatus(http.StatusNotFound).
				Send()
			return
		}
		singleton.Services[dat.Id].Heartbeat(dat.Attributes)
	}

	restful.Writer(w).
		WithResourceArr(heartbeatResources).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()
}

func HanlderFuncGetClusterServices(w http.ResponseWriter, r *http.Request) {

	var (
		vars     = mux.Vars(r)
		id       = vars["id"]
		ok       = false
		list     = map[string]struct{}{}
		services = map[string]*Service{}
	)

	if list, ok = singleton.Clusters[id]; !ok {
		restful.Writer(w).
			WithResources(
				restful.NewErrorResource(
					restful.Error404NotFound(),
				),
			).
			WithLinks(
				&resources.Links{
					Self: r.URL.String(),
				},
			).
			WithHeader(
				constants.HeaderContentTypeJSON.Key,
				constants.HeaderContentTypeJSON.Value,
			).
			WithStatus(http.StatusNotFound).
			Send()
		return
	}

	for key, _ := range list {
		services[key] = singleton.Services[key]
	}

	restful.Writer(w).
		WithResources(
			NewServicesResource(
				&services,
			),
		).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()
}

func HanlderFuncPostClusterServices(w http.ResponseWriter, r *http.Request) {

	var (
		vars              = mux.Vars(r)
		id                = vars["id"]
		ok                = false
		services          = map[string]*Service{}
		servicesResources = new([]ServiceResource)
		exists            = struct{}{}
		list              map[string]struct{}
	)

	if list, ok = singleton.Clusters[id]; !ok {
		// initialise if not exists
		list = map[string]struct{}{}
	}

	restful.Reader(r).
		WithResourceArr(servicesResources).
		Receive()

	for _, val := range *servicesResources {
		if _, ok := singleton.Services[val.Id]; !ok {
			restful.Writer(w).
				WithResources(
					restful.NewErrorResource(
						restful.Error404NotFound(),
					),
				).
				WithLinks(
					&resources.Links{
						Self: r.URL.String(),
					},
				).
				WithHeader(
					constants.HeaderContentTypeJSON.Key,
					constants.HeaderContentTypeJSON.Value,
				).
				WithStatus(http.StatusNotFound).
				Send()
			return
		}
		// mark service to cluster
		list[val.Id] = exists
		services[val.Id] = singleton.Services[val.Id]
	}

	singleton.Clusters[id] = list

	// return cluster services
	restful.Writer(w).
		WithResources(
			NewServicesResource(
				&services,
			),
		).
		WithLinks(
			&resources.Links{
				Self: r.URL.String(),
			},
		).
		WithHeader(
			constants.HeaderContentTypeJSON.Key,
			constants.HeaderContentTypeJSON.Value,
		).
		WithStatus(http.StatusOK).
		Send()
}
