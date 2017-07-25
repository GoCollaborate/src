package cmd

import (
	"flag"
	"github.com/GoCollaborate/constants"
)

type SysVars struct {
	ServerMode      string
	DebugMode       bool
	Port            int
	ContactBookPath string
	LogPath         string
	DataStorePath   string
	MaxRoutines     int
}

func InitCmdEnv() *SysVars {
	svrMode := flag.String("svrmode", constants.CollaboratorMode, "define the server mode, clbt=collaborator, cdnt=coordinator")
	debug := flag.Bool("debug", constants.DebugInactivated, "define if the server should be running in debug mode")
	port := flag.Int("port", constants.DefaultListenPort, "define the port the server should listen to")
	ctBookPath := flag.String("ctb", constants.DefaultContactBookPath, "define the path to store remote networking contacts")
	logPath := flag.String("log", constants.DefaultLogPath, "define the path to store GoCollaborate log files")
	dsPath := flag.String("ds", constants.DefaultDataStorePath, "define the path to store GoCollaborate registry data files")
	mxRoutines := flag.Int("mxrt", constants.DefaultWorkerPerMaster, "define the maximum number of threads working simultaneously to consume tasks, we recommend the number not to be greater than 1000")
	flag.Parse()
	return DefaultSysVars(&SysVars{*svrMode, *debug, *port, *ctBookPath, *logPath, *dsPath, *mxRoutines})
}

func DefaultSysVars(sysvars *SysVars) *SysVars {
	// fill out empty default system variables here
	return sysvars
}

// any change of sysvars should be processed through this function call
func (sysVars *SysVars) SetVar(func(v interface{})) {

}
