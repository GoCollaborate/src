package constants

import (
	"errors"
	"os"
	"time"
)

// system vars
const (
	CollaboratorModeAbbr   = "clbt"
	CollaboratorMode       = "collaborator"
	CoordinatorModeAbbr    = "cdnt"
	CoordinatorMode        = "coordinator"
	DebugInactivated       = false
	DebugActivated         = true
	DefaultListenPort      = 8080
	DefaultCasePath        = "./case.json"
	DefaultLogPath         = "./history.log"
	DefaultDataStorePath   = "./collaborate.dat"
	DefaultLogPrefix       = ""
	CleanHistory           = true
	DefaultNotCleanHistory = false
	Mapper                 = "Mapper"
	Reducer                = "Reducer"
	Executor               = "Executor"
	Function               = "Function"
	HashFunction           = "HashFunction"
	Shared                 = "Shared"
	Local                  = "Local"
	Limit                  = "Limit"
	ProjectPath            = "ProjectPath"
	ProjectUnixPath        = "ProjectUnixPath"
	DefaultWorkerPerMaster = 10
	DefaultHost            = "localhost"
	DefaultGossipNum       = 5
	DefaultCaseID          = "GoCollaborateStandardCase"
	DefaultJobRequestBurst = 1000
)

// store setting
const (
	DefaultHashLength = 12
)

var (
	DefaultReadTimeout              = 15 * time.Second
	DefaultPeriodShort              = 500 * time.Millisecond
	DefaultPeriodLong               = 2000 * time.Millisecond
	DefaultPeriodRoutineDay         = 24 * time.Hour
	DefaultPeriodRoutineWeek        = 7 * 24 * time.Hour
	DefaultPeriodRoutine30Days      = 30 * DefaultPeriodRoutineDay
	DefaultPeriodPermanent          = 0 * time.Second
	DefaultTaskExpireTime           = 30 * time.Second
	DefaultGCInterval               = 30 * time.Second
	DefaultMaxMappingTime           = 600 * time.Second
	DefaultSynInterval              = 3 * time.Minute
	DefaultJobRequestRefillInterval = 1 * time.Millisecond
)

// executor types
const (
	ExecutorTypeMapper   = "mapper"
	ExecutorTypeReducer  = "reducer"
	ExecutorTypeShuffler = "shuffler"
	ExecutorTypeCombiner = "combiner"
	ExecutorTypeDefault  = "default"
)

// communication types
const (
	ArgTypeInteger              = "integer"
	ArgTypeNumber               = "number"
	ArgTypeString               = "string"
	ArgTypeObject               = "object"
	ArgTypeBoolean              = "boolean"
	ArgTypeNull                 = "null"
	ArgTypeArray                = "array"
	ConstraintTypeMax           = "maximum"
	ConstraintTypeMin           = "minimum"
	ConstraintTypeXMin          = "exclusiveMinimum"
	ConstraintTypeXMax          = "exclusiveMaximum"
	ConstraintTypeUniqueItems   = "uniqueItems"
	ConstraintTypeMaxProperties = "maxProperties"
	ConstraintTypeMinProperties = "minProperties"
	ConstraintTypeMaxLength     = "maxLength"
	ConstraintTypeMinLength     = "minLength"
	ConstraintTypePattern       = "pattern"
	ConstraintTypeMaxItems      = "maxItems"
	ConstraintTypeMinItems      = "minItems"
	ConstraintTypeEnum          = "enum" // value of interface{} should accept a slice
	ConstraintTypeAllOf         = "allOf"
	ConstraintTypeAnyOf         = "anyOf"
	ConstraintTypeOneOf         = "oneOf"
)

// errors
var (
	ErrUnknownCmdArg                   = errors.New("GoCollaborate: unknown commandline argument, please enter -h to check out")
	ErrConnectionClosed                = errors.New("GoCollaborate: connection closed")
	ErrUnknown                         = errors.New("GoCollaborate: unknown error")
	ErrAPIError                        = errors.New("GoCollaborate: api error")
	ErrNoCollaborator                  = errors.New("GoCollaborate: collaborator does not exist")
	ErrCollaboratorExists              = errors.New("GoCollaborate: collaborator already exists")
	ErrNoService                       = errors.New("GoCollaborate: service of id does not exist")
	ErrConflictService                 = errors.New("GoCollaborate: found conflict, service of id already exists")
	ErrNoRegister                      = errors.New("GoCollaborate: register does not exist")
	ErrConflictRegister                = errors.New("GoCollaborate: found conflict, provider of the service already exists")
	ErrNoSubscriber                    = errors.New("GoCollaborate: subscriber does not exist")
	ErrConflictSubscriber              = errors.New("GoCollaborate: found conflict, subscriber of the service already exists")
	ErrTimeout                         = errors.New("GoCollaborate: task timeout error")
	ErrNoPeers                         = errors.New("GoCollaborate: no peer appears in the contact book")
	ErrFunctNotExist                   = errors.New("GoCollaborate: no such function found in store")
	ErrJobNotExist                     = errors.New("GoCollaborate: no sucn job found in store")
	ErrExecutorNotFound                = errors.New("GoCollaborate: no such executor found in store")
	ErrLimiterNotFound                 = errors.New("GoCollaborate: no such limiter found in store")
	ErrValNotFound                     = errors.New("GoCollaborate: no value found with such key")
	ErrCaseMismatch                    = errors.New("GoCollaborate: case mismatch error")
	ErrCollaboratorMismatch            = errors.New("GoCollaborate: collaborator mismatch error")
	ErrUnknownMsgType                  = errors.New("GoCollaborate: unknown message type error")
	ErrMapTaskFailing                  = errors.New("GoCollaborate: map operation failing error")
	ErrReduceTaskFailing               = errors.New("GoCollaborate: reduce operation failing error")
	ErrExecutorStackLengthInconsistent = errors.New("GoCollaborate: executor stack length inconsistent error")
	ErrMessageChannelDirty             = errors.New("GoCollaborate: message channel has unconsumed message error")
	ErrTaskChannelDirty                = errors.New("GoCollaborate: task channel has unconsumed task error")
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HTTP headers
var (
	Header200OK               = Header{"200", "OK"}
	Header201Created          = Header{"201", "Created"}
	Header202Accepted         = Header{"202", "Accepted"}
	Header204NoContent        = Header{"204", "NoContent"}
	Header403Forbidden        = Header{"403", "Forbidden"}
	Header404NotFound         = Header{"404", "NotFound"}
	Header409Conflict         = Header{"409", "Conflict"}
	Header422ExceedLimit      = Header{"422", "ExceedLimit"}
	HeaderContentTypeJSON     = Header{"Content-Type", "application/json"}
	HeaderContentTypeText     = Header{"Content-Type", "text/html"}
	HeaderCORSEnableAllOrigin = Header{"Access-Control-Allow-Origin", "*"}
)

// Gossip Protocol headers
var (
	GossipHeaderOK                   = Header{"200", "OK"}
	GossipHeaderUnknownMsgType       = Header{"400", "UnknownMessageType"}
	GossipHeaderCaseMismatch         = Header{"401", "CaseMismatch"}
	GossipHeaderCollaboratorMismatch = Header{"401", "CollaboratorMismatch"}
	GossipHeaderUnknownError         = Header{"500", "UnknownGossipError"}
)

// restful
const (
	JSONAPIVersion = `{"version":"1.0"}`
)

var (
	ProjectDir     = ""
	ProjectUnixDir = ""
	LibDir         = "github.com/GoCollaborate/"
	LibUnixDir     = os.Getenv("GOPATH") + "/src/github.com/GoCollaborate/"
)
