package main

import (
	"github.com/GoCollaborate/cmd"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote"
	"github.com/GoCollaborate/server"
	"github.com/GoCollaborate/utils"
	"github.com/gorilla/mux"
	"net/http"
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

	// initialise server
	router := mux.NewRouter()
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
		utils.AdaptRouterToDebugMode(router)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	pbls.Distribute()
}
