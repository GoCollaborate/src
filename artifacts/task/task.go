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
		return constants.DEFAULT_PERIOD_SHORT
	case LONG:
		return constants.DEFAULT_PERIOD_LONG
	case PERMANENT:
		return constants.DEFAULT_PERIOD_PERMANENT
	default:
		return constants.DEFAULT_PERIOD_PERMANENT
	}
}

func (t *taskPriority) GetPriority() taskPriority {
	return *t
}

func NewCollection() *Collection {
	return &Collection{}
}

func (cg *Collection) Append(cs ...interface{}) *Collection {
	*cg = append(*cg, cs...)
	return cg
}

func (cg *Collection) IsEmpty() bool {
	return len(*cg) == 0
}

func (cg *Collection) Length() int {
	return len(*cg)
}

func (cg *Collection) Filter(f func(interface{}) bool) *Collection {
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

type Collection []interface{}
