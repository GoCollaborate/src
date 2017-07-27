package utils

import (
	"github.com/GoCollaborate/constants"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
)

func AdaptRouterToDebugMode(router *mux.Router) *mux.Router {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	return router
}

func AdaptHTTPWithHeader(w http.ResponseWriter,
	header constants.Header) http.ResponseWriter {
	w.Header().Add(header.Key, header.Value)
	return w
}
