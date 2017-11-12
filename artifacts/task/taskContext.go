package task

import (
	"github.com/GoCollaborate/src/constants"
	"github.com/GoCollaborate/src/utils"
	"sync"
)

var lock sync.RWMutex = sync.RWMutex{}

type TaskContext struct {
	Context map[string]interface{}
}

func NewTaskContext(ctx interface{}) *TaskContext {
	maps := utils.Map(ctx)
	return &TaskContext{maps}
}

func (this *TaskContext) Entries() map[string]interface{} {
	return this.Context
}

func (this *TaskContext) Set(key string, val interface{}) {
	lock.Lock()
	this.Context[key] = val
	lock.Unlock()
}

func (this *TaskContext) Get(key string) (interface{}, error) {
	if val := this.Context[key]; val != nil {
		return val, nil
	}
	return nil, constants.ErrValNotFound
}
