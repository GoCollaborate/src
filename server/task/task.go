package task

import (
	"github.com/GoCollaborate/constants"
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

type Task struct {
	Type       taskType
	Priority   taskPriority
	Consumable func([]Countable, []Countable, *TaskContext) bool
	Source     []Countable
	Result     []Countable
	Context    *TaskContext
}

type Countable interface {
	Int() int
	Int64() int64
	Float64() float64
	Bool() bool
	String() string
}

type Chunk struct {
	Value interface{}
}

func (c *Chunk) Int() int {
	return c.Value.(int)
}

func (c *Chunk) Int64() int64 {
	return c.Value.(int64)
}

func (c *Chunk) Float64() float64 {
	return c.Value.(float64)
}

func (c *Chunk) Bool() bool {
	return c.Value.(bool)
}

func (c *Chunk) String() string {
	return c.Value.(string)
}

type Split struct {
	Value interface{}
}

func (s *Split) Int() int {
	return s.Value.(int)
}

func (s *Split) Int64() int64 {
	return s.Value.(int64)
}

func (s *Split) Float64() float64 {
	return s.Value.(float64)
}

func (s *Split) Bool() bool {
	return s.Value.(bool)
}

func (s *Split) String() string {
	return s.Value.(string)
}
