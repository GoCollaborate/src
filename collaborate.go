package main

import (
	"github.com/GoCollaborate/cmd"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote"
	"github.com/GoCollaborate/server"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"time"
)

var pbls *server.Publisher

func main() {
	// initialise environment
	sysvars := cmd.InitCmdEnv()
	var (
		localLogger *logger.Logger
		logFile     *os.File
	)

	switch sysvars.CleanHistory {
	case constants.CleanHistory:
		localLogger, logFile = logger.NewLogger(sysvars.LogPath, constants.DefaultLogPrefix, true)
	default:
		localLogger, logFile = logger.NewLogger(sysvars.LogPath, constants.DefaultLogPrefix, true)
	}
	// set handler for router
	router := mux.NewRouter()
	switch sysvars.DebugMode {
	case true:
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
	default:
	}

	switch sysvars.ServerMode {
	case constants.CollaboratorModeAbbr, constants.CollaboratorMode:
		// create book keeper
		bkp := new(remote.BookKeeper)
		// create contact book
		contactBook := remote.ContactBook{[]remote.Agent{}, remote.Agent{}, *remote.Default(), false, false, time.Now().Unix()}
		// lock book keeper to contact book
		bkp.LookAtAndWatch(&contactBook)

		// pbls = server.NewPublisher(localLogger)
		mst := server.NewMaster(localLogger)
		// pbls.Connect(mst)
		mst.BatchAttach(sysvars.MaxRoutines)
		mst.LaunchAll()
		bkp.Handle(router)
	case constants.CoordinatorModeAbbr, constants.CoordinatorMode:
		regCenter := remote.GetRegCenterInstance(sysvars.Port, localLogger)
		regCenter.Handle(router)
	}

	// launch server
	serv := &http.Server{
		Addr:        ":" + strconv.Itoa(sysvars.Port),
		Handler:     router,
		ReadTimeout: constants.DefaultReadTimeout,
	}
	err := serv.ListenAndServe()
	localLogger.LogError(err.Error())
	logFile.Close()
}
