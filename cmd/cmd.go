package cmd

import (
	"flag"
	"github.com/GoCollaborate/constants"
	"path/filepath"
	"sync"
)

type SysVars struct {
	ServerMode      string
	DebugMode       bool
	Port            int
	CasePath        string
	LogPath         string
	DataStorePath   string
	MaxRoutines     int
	CleanHistory    bool
	WorkerPerMaster int
	GossipNum       int
	CaseID          string
}

var singleton *SysVars
var once sync.Once

func InitCmdEnv() *SysVars {
	svrMode := flag.String("svrmode", constants.CollaboratorMode, "define the server mode, clbt=collaborator, cdnt=coordinator")
	debug := flag.Bool("debug", constants.DebugInactivated, "define if the server should be running in debug mode")
	port := flag.Int("port", constants.DefaultListenPort, "define the port the server should listen to")
	casePath := flag.String("cpath", constants.DefaultCasePath, "define the path to store CASE information")
	logPath := flag.String("log", constants.DefaultLogPath, "define the path to store GoCollaborate log files")
	dsPath := flag.String("ds", constants.DefaultDataStorePath, "define the path to store GoCollaborate data files")
	mxRoutines := flag.Int("mxrt", constants.DefaultWorkerPerMaster, "define the maximum number of routines working simultaneously to consume TASK, we recommend the number not to be greater than 1000")
	cleanHis := flag.Bool("clean", constants.DefaultNotCleanHistory, "define if the program should clean previous history while launching up (note: the cleaning is irrecoverable!)")
	numWorkers := flag.Int("numwks", constants.DefaultWorkerPerMaster, "define the number of WORKER created for each MASTER)")
	numGossip := flag.Int("numgossip", constants.DefaultGossipNum, "define the number of nodes to spread to while exchanging cards)")
	caseID := flag.String("cid", constants.DefaultCaseID, "define the identifier of CASE cluster, only message within the same CASE cluster will be processed by each COLLABORATOR")
	flag.Parse()
	if !flag.Parsed() {
		panic(constants.ErrUnknownCmdArg)
	}
	Init()
	return &SysVars{*svrMode, *debug, *port, *casePath, *logPath, *dsPath, *mxRoutines, *cleanHis, *numWorkers, *numGossip, *caseID}
}

func Combine(vars ...*SysVars) *SysVars {
	if len(vars) < 1 {
		return singleton
	}
	for _, s := range vars {
		singleton.ServerMode = s.ServerMode
		singleton.DebugMode = s.DebugMode
		singleton.Port = s.Port
		singleton.CasePath = s.CasePath
		singleton.LogPath = s.LogPath
		singleton.DataStorePath = s.DataStorePath
		singleton.MaxRoutines = s.MaxRoutines
		singleton.CleanHistory = s.CleanHistory
		singleton.WorkerPerMaster = s.WorkerPerMaster
		singleton.GossipNum = s.GossipNum
		singleton.CaseID = s.CaseID
	}
	cp := *singleton
	return &cp
}

func Vars() *SysVars {
	cp := *singleton
	return &cp
}

func Init() {
	once.Do(func() {
		singleton = &SysVars{}
		path, err := filepath.Abs("./")
		if err != nil {
			panic(err)
		}
		constants.ProjectDir = path
	})
}
