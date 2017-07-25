package main

import (
	"github.com/GoCollaborate/cmd"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logging"
	"github.com/GoCollaborate/remote"
	"github.com/GoCollaborate/server"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"
)

var pbls *server.Publisher

func main() {
	// initialise environment
	sysvars := cmd.InitCmdEnv()
	localLogger, logFile := logger.NewLogger(sysvars.LogPath, constants.DefaultLogPrefix)

	// initialise server
	router := mux.NewRouter()
	switch sysvars.ServerMode {
	case constants.CollaboratorModeAbbr, constants.CollaboratorMode:
		// create contact book
		contactBook := remote.ContactBook{[]remote.Agent{}, remote.Agent{}, *remote.Default(), false, false, time.Now().Unix()}
		contactBook.Load()
		contactBook.LaunchServer()

		// bridge to remote servers
		contactBook.Bridge()

		/*
			terminate the rpc connection at any time by calling Bridge() function
		*/
		// contactBook.Terminate()

		pbls = server.NewPublisher(localLogger)
		mst := server.NewMaster(localLogger)
		pbls.Connect(mst)
		mst.BatchAttach(sysvars.MaxRoutines)
		mst.LaunchAll()
		router.HandleFunc("/", index)
		serv := &http.Server{
			Addr:        ":" + strconv.Itoa(sysvars.Port),
			Handler:     router,
			ReadTimeout: constants.DefaultReadTimeout,
		}
		err := serv.ListenAndServe()
		localLogger.LogError(err.Error())
		logFile.Close()
	case constants.CoordinatorModeAbbr, constants.CoordinatorMode:
		regCenter := remote.NewRegCenter(sysvars.Port)
		regCenter.Listen()
	}
	if sysvars.DebugMode {
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	pbls.Distribute()
}
