package main

import (
	"github.com/GoCollaborate/logging"
	//"github.com/GoCollaborate/remote"
	"github.com/GoCollaborate/server"
	"github.com/gorilla/mux"
	"net/http"
	// "fmt"
	"net/http/pprof"
	"strconv"
	"time"
)

var pbls *server.Publisher

func main() {
	// config := remote.Config{[]remote.Agent{}, time.Now().Unix(), *remote.Default()}
	// config.Load()
	// config.LaunchServer()

	// config.Bridge()
	// fmt.Println("Update...")
	// Delay(10)
	// fmt.Println("EndDelay...")
	// config.Terminate()
	localLogger, logFile := logger.NewLogger("./history.log", "GoCollaborate/")
	pbls = server.NewPublisher(localLogger)
	mst := server.NewMaster(localLogger)
	pbls.Connect(mst)
	mst.BatchAttach(10)
	mst.LaunchAll()
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))

	serv := &http.Server{
		Addr:        ":" + strconv.Itoa(8080),
		Handler:     router,
		ReadTimeout: 15 * time.Second,
	}
	serv.ListenAndServe()
	for {
	}
	logFile.Close()
}

func index(w http.ResponseWriter, r *http.Request) {
	pbls.Distribute()
}

// func Delay(sec time.Duration) {
// 	tm := time.NewTimer(sec * time.Second)
// 	<-tm.C
// }
