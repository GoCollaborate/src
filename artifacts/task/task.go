package task

import (
	"github.com/GoCollaborate/src/constants"
	"time"
)

type taskType int

const (
	SHORT taskType = iota
	LONG
	ROUTINE
	PERMANENT
)

type taskPriority int

const (
	BASE taskPriority = iota
	LOW
	MEDIUM
	HIGH
	URGENT
)

type TaskType interface {
	GetType() taskType
	GetTimeout() time.Time
}

type TaskPriority interface {
	GetPriority() taskPriority
}

func (t *taskType) GetType() taskType {
	return *t
}

// if return nil, this taks is identified as an routine task
func (t *taskType) GetTimeout() time.Duration {
	switch t.GetType() {
	case SHORT:
		return constants.DefaultPeriodShort
	case LONG:
		return constants.DefaultPeriodLong
	case PERMANENT:
		return constants.DefaultPeriodPermanent
	default:
		return constants.DefaultPeriodPermanent
	}
}

func (t *taskPriority) GetPriority() taskPriority {
	return *t
}

func (cg *Collection) Append(cs ...Countable) *Collection {
	*cg = append(*cg, cs...)
	return cg
}

func (cg *Collection) IsEmpty() bool {
	return len(*cg) == 0
}

func (cg *Collection) Length() int {
	return len(*cg)
}

func (cg *Collection) Filter(f func(Countable) bool) *Collection {
	var (
		clct = Collection{}
	)

	for _, c := range *cg {
		if f(c) {
			clct = append(clct, c)
		}
	}

	*cg = clct

	return cg
}

type Task struct {
	Type       taskType
	Priority   taskPriority
	Consumable string
	Source     Collection
	Result     Collection
	Context    *TaskContext
	Stage      int
}

type Collection []Countable

type Countable interface{}
