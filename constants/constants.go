package constants

import (
	"time"
)

const (
	DefaultWorkerPerMaster    = 10
	DefaultHost               = "localhost"
	DefaultPeriodShort        = 500 * time.Millisecond
	DefaultPeriodLong         = 2000 * time.Millisecond
	DefaultPeriodRoutineDay   = 24 * time.Hour
	DefaultPeriodRoutineWeek  = 7 * 24 * time.Hour
	DefaultPeriodRoutineMonth = 1 * time.Month
)

var (
	ErrConnectionClosed   = errors.New("GoCollaborate: connection closed")
	ErrUnknown            = errors.New("GoCollaborate: unknown error")
	ErrAPIError           = errors.New("GoCollaborate: api error")
	ErrNoCollaborator     = errors.New("GoCollaborate: collaborator does not exist")
	ErrCollaboratorExists = errors.New("GoCollaborate: collaborator already exists")
)
