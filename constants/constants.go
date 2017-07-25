package constants

import (
	"errors"
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
	DefaultLogPrefix       = "GoCollaborate/"
)

// master/worker setting
const (
	DefaultWorkerPerMaster = 10
	DefaultHost            = "localhost"
)

var (
	DefaultReadTimeout         = 15 * time.Second
	DefaultPeriodShort         = 500 * time.Millisecond
	DefaultPeriodLong          = 2000 * time.Millisecond
	DefaultPeriodRoutineDay    = 24 * time.Hour
	DefaultPeriodRoutineWeek   = 7 * 24 * time.Hour
	DefaultPeriodRoutine30Days = 30 * DefaultPeriodRoutineDay
	DefaultPeriodPermanent     = 0 * time.Second
)

// communication types
const (
	ArgTypeInteger          = "integer"
	ArgTypeNumber           = "number"
	ArgTypeString           = "string"
	ArgTypeObject           = "object"
	ArgTypeBoolean          = "boolean"
	ArgTypeNull             = "null"
	ArgTypeArray            = "array"
	ConstraintTypeMax       = "maximum"
	ConstraintTypeMin       = "minimum"
	ConstraintTypeMaxLength = "maxLength"
	ConstraintTypeMinLength = "minLength"
	ConstraintTypePattern   = "pattern"
	ConstraintTypeMaxItems  = "maxItems"
	ConstraintTypeMinItems  = "minItems"
	ConstraintTypeEnum      = "enum" // value of interface{} should accept a slice
)

// errors
var (
	ErrUnknownCmdArg      = errors.New("GoCollaborate: unknown commandline argument, please enter -h to check out")
	ErrConnectionClosed   = errors.New("GoCollaborate: connection closed")
	ErrUnknown            = errors.New("GoCollaborate: unknown error")
	ErrAPIError           = errors.New("GoCollaborate: api error")
	ErrNoCollaborator     = errors.New("GoCollaborate: collaborator does not exist")
	ErrCollaboratorExists = errors.New("GoCollaborate: collaborator already exists")
)
