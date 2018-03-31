package cmd

import (
	"flag"
	"github.com/GoCollaborate/src/constants"
	"path/filepath"
	"strconv"
	"sync"
)

type SysVars struct {
	DebugMode       bool
	Port            int
	CasePath        string
	LogPath         string
	CleanHistory    bool
	WorkerPerMaster int
	GossipNum       int
	CaseID          string
}

var singleton *SysVars
var once sync.Once

func InitCmdEnv() *SysVars {
	debug := flag.Bool("debug", constants.DEBUG_INACTIVATED, "define if the server should be running in debug mode")
	port := flag.Int("port", constants.DEFAULT_LISTEN_PORT, "define the port the server should listen to")
	casePath := flag.String("cpath", constants.DEFAULT_CASE_PATH, "define the path to store CASE information")
	logPath := flag.String("log", constants.DEFAULT_LOG_PATH, "define the path to store GoCollaborate log files")
	cleanHis := flag.Bool("clean", constants.DEFAULT_NOT_CLEAN_HISTORY, "define if the program should clean previous history while launching up (note: the cleaning is irrecoverable!)")
	numWorkers := flag.Int("numwks", constants.DEFAULT_WORKER_PER_MASTER, "define the number of WORKER created for each MASTER")
	numGossip := flag.Int("numgossip", constants.DEFAULT_GOSSIP_NUM, "define the number of nodes to spread to while exchanging cards")
	caseID := flag.String("cid", constants.DEFAULT_CASE_ID, "define the identifier of CASE cluster, only message within the same CASE cluster will be processed by each COLLABORATOR")
	flag.Parse()
	if !flag.Parsed() {
		panic(constants.ERR_UNKNOWN_CMD_ARG)
	}
	Init()
	return &SysVars{*debug, *port, *casePath, *logPath, *cleanHis, *numWorkers, *numGossip, *caseID}
}

func Combine(vars ...*SysVars) *SysVars {
	if len(vars) < 1 {
		return singleton
	}
	for _, s := range vars {
		singleton.DebugMode = s.DebugMode
		singleton.Port = s.Port
		singleton.CasePath = s.CasePath
		singleton.LogPath = s.LogPath
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

func VarsJSONArrayStr() string {
	return `{
		"headers":["CaseID",
					"DebugMode",
					"Port",
					"CasePath",
					"LogPath",
					"CleanHistory",
					"WorkerPerMaster",
					"GossipNum"],
		"data":[["CaseID","string","` + singleton.CaseID + `"],
				["DebugMode","bool",` + strconv.FormatBool(singleton.DebugMode) + `],
				["Port","int",` + strconv.Itoa(singleton.Port) + `],
				["CasePath","string","` + singleton.CasePath + `"],
				["LogPath","string","` + singleton.LogPath + `"],
				["CleanHistory","bool",` + strconv.FormatBool(singleton.CleanHistory) + `],
				["WorkerPerMaster","int",` + strconv.Itoa(singleton.WorkerPerMaster) + `],
				["GossipNum","int",` + strconv.Itoa(singleton.GossipNum) + `]]}`
}

func Init() {
	once.Do(func() {
		singleton = &SysVars{}
		path, err := filepath.Abs("./")
		if err != nil {
			panic(err)
		}
		constants.PROJECT_DIR = path
	})
}
