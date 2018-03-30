package constants

import (
	"errors"
	"os"
	"time"
)

// system vars
const (
	DEBUG_INACTIVATED         = false
	DEBUG_ACTIVATED           = true
	DEFAULT_LISTEN_PORT       = 8080
	DEFAULT_CASE_PATH         = "./case.json"
	DEFAULT_LOG_PATH          = "./history.log"
	DEFAULT_LOG_PREFIX        = ""
	CLEAN_HISTORY             = true
	DEFAULT_NOT_CLEAN_HISTORY = false
	MAPPER                    = "Mapper"
	REDUCER                   = "Reducer"
	EXECUTOR                  = "Executor"
	FUNCTION                  = "Function"
	HASH_FUNCTION             = "HashFunction"
	SHARED                    = "Shared"
	LOCAL                     = "Local"
	LIMIT                     = "Limit"
	PROJECT_PATH              = "ProjectPath"
	PROJECT_UNIX_PATH         = "ProjectUnixPath"
	DEFAULT_WORKER_PER_MASTER = 5
	DEFAULT_HOST              = "localhost"
	DEFAULT_GOSSIP_NUM        = 5
	DEFAULT_CASE_ID           = "GoCollaborateStandardCase"
	DEFAULT_JOB_REQUEST_BURST = 1000
)

// store setting
const (
	DEFAULT_HASH_LENGTH = 12
)

// time consts
var (
	DEFAULT_RPC_DIAL_TIMEOUT             = 20 * time.Second
	DEFAULT_READ_TIMEOUT                 = 15 * time.Second
	DEFAULT_PERIOD_SHORT                 = 500 * time.Millisecond
	DEFAULT_PERIOD_LONG                  = 2000 * time.Millisecond
	DEFAULT_PERIOD_ROUTINE_DAY           = 24 * time.Hour
	DEFAULT_PERIOD_ROUTINE_WEEK          = 7 * 24 * time.Hour
	DEFAULT_PERIOD_ROUTINE_30_DAYS       = 30 * DEFAULT_PERIOD_ROUTINE_DAY
	DEFAULT_PERIOD_PERMANENT             = 0 * time.Second
	DEFAULT_TASK_EXPIRY_TIME             = 30 * time.Second
	DEFAULT_GC_INTERVAL                  = 30 * time.Second
	DEFAULT_MAX_MAPPING_TIME             = 600 * time.Second
	DEFAULT_SYNC_INTERVAL                = 3 * time.Minute
	DEFAULT_HEARTBEAT_INTERVAL           = 5 * time.Second
	DEFAULT_JOB_REQUEST_REFILL_INTERVAL  = 1 * time.Millisecond
	DEFAULT_STAT_FLUSH_INTERVAL          = 20 * time.Millisecond
	DEFAULT_STAT_ABSTRACT_INTERVAL       = 3 * time.Second
	DEFAULT_COLLABORATOR_EXPIRY_INTERVAL = 10 * time.Minute
)

// executor types
const (
	EXECUTOR_TYPE_MAPPER   = "mapper"
	EXECUTOR_TYPE_REDUCER  = "reducer"
	EXECUTOR_TYPE_COMBINER = "combiner"
	EXECUTOR_TYPE_DEFAULT  = "default"
)

// communication types
const (
	ARG_TYPE_INTEGER               = "integer"
	ARG_TYPE_NUMBER                = "number"
	ARG_TYPE_STRING                = "string"
	ARG_TYPE_OBJECT                = "object"
	ARG_TYPE_BOOLEAN               = "boolean"
	ARG_TYPE_NULL                  = "null"
	ARG_TYPE_ARRAY                 = "array"
	CONSTRAINT_TYPE_MAX            = "maximum"
	CONSTRAINT_TYPE_MIN            = "minimum"
	CONSTRAINT_TYPE_XMAX           = "exclusiveMaximum"
	CONSTRAINT_TYPE_XMIN           = "exclusiveMinimum"
	CONSTRAINT_TYPE_UNIQUE_ITEMS   = "uniqueItems"
	CONSTRAINT_TYPE_MAX_PROPERTIES = "maxProperties"
	CONSTRAINT_TYPE_MIN_PROPERTIES = "minProperties"
	CONSTRAINT_TYPE_MAX_LENGTH     = "maxLength"
	CONSTRAINT_TYPE_MIN_LENGTH     = "minLength"
	CONSTRAINT_TYPE_PATTERN        = "pattern"
	CONSTRAINT_TYPE_MAX_ITEMS      = "maxItems"
	CONSTRAINT_TYPE_MIN_ITEMS      = "minItems"
	CONSTRAINT_TYPE_ENUM           = "enum" // value of interface{} should accept a slice
	CONSTRAINT_TYPE_ALLOF          = "allOf"
	CONSTRAINT_TYPE_ANYOF          = "anyOf"
	CONSTRAINT_TYPE_ONEOF          = "oneOf"
)

const (
	STATS_POLICY_SUM_OF_INTS = "StatsPolicySumOfInt"
)

// errors
var (
	ERR_UNKNOWN_CMD_ARG                    = errors.New("GoCollaborate: unknown commandline argument, please enter -h to check out")
	ERR_CONNECTION_CLOSED                  = errors.New("GoCollaborate: connection closed")
	ERR_UNKNOWN                            = errors.New("GoCollaborate: unknown error")
	ERR_API                                = errors.New("GoCollaborate: api error")
	ERR_NO_COLLABORATOR                    = errors.New("GoCollaborate: collaborator does not exist")
	ERR_COLLAOBRATOR_ALREADY_EXISTS        = errors.New("GoCollaborate: collaborator already exists")
	ERR_TASK_TIMEOUT                       = errors.New("GoCollaborate: task timeout error")
	ERR_NO_PEERS                           = errors.New("GoCollaborate: no peer appears in the contact")
	ERR_FUNCT_NOT_EXIST                    = errors.New("GoCollaborate: no such function found in store")
	ERR_JOB_NOT_EXIST                      = errors.New("GoCollaborate: no sucn job found in store")
	ERR_EXECUTOR_NOT_FOUND                 = errors.New("GoCollaborate: no such executor found in store")
	ERR_LIMITER_NOT_FOUND                  = errors.New("GoCollaborate: no such limiter found in store")
	ERR_VAL_NOT_FOUND                      = errors.New("GoCollaborate: no value found with such key")
	ERR_CASE_MISMATCH                      = errors.New("GoCollaborate: case mismatch error")
	ERR_COLLABORATOR_MISMATCH              = errors.New("GoCollaborate: collaborator mismatch error")
	ERR_UNKNOWN_MSG_TYPE                   = errors.New("GoCollaborate: unknown message type error")
	ERR_TASK_MAP_FAIL                      = errors.New("GoCollaborate: map operation failing error")
	ERR_TASK_REDUCE_FAIL                   = errors.New("GoCollaborate: reduce operation failing error")
	ERR_EXECUTOR_STACK_LENGTH_INCONSISTENT = errors.New("GoCollaborate: executor stack length inconsistent error")
	ERR_MSG_CHAN_DIRTY                     = errors.New("GoCollaborate: message channel has unconsumed message error")
	ERR_TASK_CHAN_DIRTY                    = errors.New("GoCollaborate: task channel has unconsumed task error")
	ERR_STAT_TYPE_NOT_FOUND                = errors.New("GoCollaborate: stat type not found error")
	ERR_INPUT_STREAM_CORRUPTED             = errors.New("GoCollaborate: input stream corrupted error")
	ERR_INPUT_STREAM_NOT_SUPPORTED         = errors.New("GoCollaborate: input stream type not suppoted error")
	ERR_IO_DECODE_POINTER_REQUIRED         = errors.New("GoCollaborate: Decode error, the reference instance must be a pointer")
	ERR_IO_DECODE_SLICE_REQUIRED           = errors.New("GoCollaborate: Decode error, the reference instance must be a slice")
	ERR_IO_DECODE_STRUCT_REQUIRED          = errors.New("GoCollaborate: Decode error, the reference instance must be a struct")
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HTTP headers
var (
	HEADER_CONTENT_TYPE_JSON       = Header{"Content-Type", "application/json"}
	HEADER_CONTENT_TYPE_TEXT       = Header{"Content-Type", "text/html"}
	HEADER_CORS_ENABLED_ALL_ORIGIN = Header{"Access-Control-Allow-Origin", "*"}
)

// Gossip Protocol headers
var (
	GOSSIP_HEADER_OK                    = Header{"200", "OK"}
	GOSSIP_HEADER_UNKNOWN_MSG_TYPE      = Header{"400", "UnknownMessageType"}
	GOSSIP_HEADER_CASE_MISMATCH         = Header{"401", "CaseMismatch"}
	GOSSIP_HEADER_COLLABORATOR_MISMATCH = Header{"401", "CollaboratorMismatch"}
	GOSSIP_HEADER_UNKNOWN_ERROR         = Header{"500", "UnknownGossipError"}
)

// restful
const (
	JSON_API_VERSION = `{"version":"1.0"}`
)

var (
	PROJECT_DIR      = ""
	PROJECT_UNIX_DIR = ""
	LIB_DIR          = "github.com/GoCollaborate/src/"
	LIB_UNIX_DIR     = os.Getenv("GOPATH") + "/src/github.com/GoCollaborate/src/"
)
