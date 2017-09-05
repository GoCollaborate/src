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
	DefaultContactBookPath = "./contact.json"
	DefaultLogPath         = "./history.log"
	DefaultDataStorePath   = "./collaborate.dat"
	DefaultLogPrefix       = "GoCollaborate:"
	CleanHistory           = true
	DefaultNotCleanHistory = false
	Mapper                 = "Mapper"
	Reducer                = "Reducer"
	Function               = "Function"
	HashFunction           = "HashFunction"
	Shared                 = "Shared"
	Local                  = "Local"
	ProjectPath            = "ProjectPath"
)

// master/worker setting
const (
	DefaultWorkerPerMaster = 10
	DefaultHost            = "localhost"
)

// funcstore setting
const (
	DefaultHashLength = 12
)

var (
	DefaultReadTimeout         = 15 * time.Second
	DefaultPeriodShort         = 500 * time.Millisecond
	DefaultPeriodLong          = 2000 * time.Millisecond
	DefaultPeriodRoutineDay    = 24 * time.Hour
	DefaultPeriodRoutineWeek   = 7 * 24 * time.Hour
	DefaultPeriodRoutine30Days = 30 * DefaultPeriodRoutineDay
	DefaultPeriodPermanent     = 0 * time.Second
	DefaultTaskExpireTime      = 30 * time.Second
	DefaultGCInterval          = 30 * time.Second
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
	ErrUnknownCmdArg      = errors.New("GoCollaborate: unknown commandline argument, please enter -h to check out")
	ErrConnectionClosed   = errors.New("GoCollaborate: connection closed")
	ErrUnknown            = errors.New("GoCollaborate: unknown error")
	ErrAPIError           = errors.New("GoCollaborate: api error")
	ErrNoCollaborator     = errors.New("GoCollaborate: collaborator does not exist")
	ErrCollaboratorExists = errors.New("GoCollaborate: collaborator already exists")
	ErrNoService          = errors.New("GoCollaborate: service of id does not exist")
	ErrConflictService    = errors.New("GoCollaborate: found conflict, service of id already exists")
	ErrNoRegister         = errors.New("GoCollaborate: register does not exist")
	ErrConflictRegister   = errors.New("GoCollaborate: found conflict, provider of the service already exists")
	ErrNoSubscriber       = errors.New("GoCollaborate: subscriber does not exist")
	ErrConflictSubscriber = errors.New("GoCollaborate: found conflict, subscriber of the service already exists")
	ErrTimeout            = errors.New("GoCollaborate: task timeout error")
	ErrNoPeers            = errors.New("GoCollaborate: no peer appears in the contact book")
	ErrFunctNotExist      = errors.New("GoCollaborate: no such function found in store")
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// header
var (
	Header200OK        = Header{"200", "OK"}
	Header201Created   = Header{"201", "Created"}
	Header202Accepted  = Header{"202", "Accepted"}
	Header204NoContent = Header{"204", "NoContent"}
	Header403Forbidden = Header{"403", "Forbidden"}
	Header404NotFound  = Header{"404", "NotFound"}
	Header409Conflict  = Header{"409", "Conflict"}
)

// restful
const (
	JSONAPIVersion = `{"version":"1.0"}`
)

var (
	ProjectDir     = "github.com/GoCollaborate/"
	ProjectUnixDir = os.Getenv("GOPATH") + "/src/github.com/GoCollaborate/"
)
