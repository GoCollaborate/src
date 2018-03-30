package collaborate

import (
	"github.com/GoCollaborate/src/artifacts/iexecutor"
	"github.com/GoCollaborate/src/artifacts/imapper"
	"github.com/GoCollaborate/src/artifacts/ireducer"
	"github.com/GoCollaborate/src/artifacts/master"
	"github.com/GoCollaborate/src/artifacts/task"
	"github.com/GoCollaborate/src/cmd"
	"github.com/GoCollaborate/src/collaborator"
	"github.com/GoCollaborate/src/constants"
	"github.com/GoCollaborate/src/logger"
	"github.com/GoCollaborate/src/store"
	"github.com/GoCollaborate/src/utils"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Set(key string, val ...interface{}) interface{} {
	cmd.Init()
	switch key {
	case constants.Mapper:
		fs := store.GetInstance()
		fs.SetMapper(val[0].(imapper.IMapper), val[1].(string))
	case constants.Reducer:
		fs := store.GetInstance()
		fs.SetReducer(val[0].(ireducer.IReducer), val[1].(string))
	case constants.Executor:
		fs := store.GetInstance()
		fs.SetExecutor(val[0].(iexecutor.IExecutor), val[1].(string))
	case constants.Function:
		// register function
		fs := store.GetInstance()
		f := val[0].(func(source *task.Collection,
			result *task.Collection,
			context *task.TaskContext) bool)
		if len(val) > 1 {
			fs.Add(f, val[1].(string))
			break
		}
		fs.Add(f)
	case constants.HashFunction:
		// register hash function
		fs := store.GetInstance()
		f := val[0].(func(source *task.Collection,
			result *task.Collection,
			context *task.TaskContext) bool)
		return fs.HAdd(f)
	case constants.Shared:
		fs := store.GetInstance()

		methods := val[0].([]string)
		handlers := make([]func(w http.ResponseWriter, r *http.Request, bg *task.Background), len(val)-1)
		for i, v := range val[1:] {
			handlers[i] = v.(func(w http.ResponseWriter, r *http.Request, bg *task.Background))
		}

		// register jobs
		fs.AddShared(methods, handlers...)
	case constants.Local:
		fs := store.GetInstance()

		methods := val[0].([]string)
		handlers := make([]func(w http.ResponseWriter, r *http.Request, bg *task.Background), len(val)-1)
		for i, v := range val[1:] {
			handlers[i] = v.(func(w http.ResponseWriter, r *http.Request, bg *task.Background))
		}

		// register jobs
		fs.AddLocal(methods, handlers...)
	case constants.Limit:
		var (
			fs    = store.GetInstance()
			limit = constants.DefaultJobRequestRefillInterval
			burst = constants.DefaultJobRequestBurst
		)
		if len(val) > 1 {
			limit = val[1].(time.Duration)
		}
		if len(val) > 2 {
			burst = val[2].(int)
		}
		fs.SetLimiter(val[0].(string), rate.Every(limit), burst)
	case constants.ProjectPath:
		constants.ProjectDir = val[0].(string)
	case constants.ProjectUnixPath:
		constants.ProjectUnixDir = val[0].(string)
	}
	return nil
}

func Run(vars ...*cmd.SysVars) {
	// initialise environment
	cmd.Combine(append(vars, cmd.InitCmdEnv())...)

	var (
		localLogger *logger.Logger
		logFile     *os.File
	)

	runVars := cmd.Vars()

	switch runVars.CleanHistory {
	case constants.CleanHistory:
		localLogger, logFile = logger.NewLogger(runVars.LogPath, constants.DefaultLogPrefix, true)
	default:
		localLogger, logFile = logger.NewLogger(runVars.LogPath, constants.DefaultLogPrefix, true)
	}

	// set handler for router
	router := store.GetRouter()
	switch runVars.DebugMode {
	case true:
		router = utils.AdaptRouterToDebugMode(router)
	default:
	}

	logger.LogLogo("", "(c) 2017 GoCollaborate", "", "Author: Hastings Yeung", "Github: https://github.com/GoCollaborate", "")

	switch runVars.ServerMode {
	// WARNING: Deprecated since GoCollaborate ver 0.5.x
	// case constants.CoordinatorModeAbbr, constants.CoordinatorMode:
	// 	cdnt := coordinator.GetCoordinatorInstance(int32(runVars.Port))
	// 	cdnt.Handle(router)

	default:
		// create collaborator
		clbt := collaborator.NewCollaborator()

		mst := master.NewMaster()
		mst.BatchAttach(runVars.MaxRoutines)
		mst.LaunchAll()

		clbt.Join(mst)

		clbt.Handle(router)
	}

	// launch server
	serv := &http.Server{
		Addr:        ":" + strconv.Itoa(runVars.Port),
		Handler:     router,
		ReadTimeout: constants.DefaultReadTimeout,
	}
	err := serv.ListenAndServe()
	logger.LogError(err.Error())
	localLogger.LogError(err.Error())
	logFile.Close()
}
