package funcstore

import (
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/utils"
	"sync"
)

var singleton *FS
var once sync.Once

func GetFSInstance() *FS {
	once.Do(func() {
		singleton = &FS{make(map[string]func(source *[]task.Countable,
			result *[]task.Countable,
			context *task.TaskContext) chan bool), make(map[string]chan bool)}
	})
	return singleton
}

type FS struct {
	Funcs map[string]func(source *[]task.Countable,
		result *[]task.Countable,
		context *task.TaskContext) chan bool
	Outbound map[string]chan bool
}

func (fs *FS) Add(f func(source *[]task.Countable,
	result *[]task.Countable,
	context *task.TaskContext) chan bool, id ...string) {
	var i string
	if len(id) < 1 {
		i = utils.StripRouteToFunctName(utils.ReflectFuncName(f))
	} else {
		i = id[0]
	}
	fs.Funcs[i] = f
	fs.Outbound[i] = make(chan bool)
}

func (fs *FS) Call(id string, source *[]task.Countable,
	result *[]task.Countable,
	context *task.TaskContext) {

	if f := fs.Funcs[id]; f != nil {
		fs.Outbound[id] <- <-f(source, result, context)
		return
	}

	logger.LogError(constants.ErrFunctNotExist)
	return
}

func (fs *FS) Listen(id string) chan bool {
	if o := fs.Outbound[id]; o != nil {
		return o
	}
	logger.LogError(constants.ErrFunctNotExist)
	out := make(chan bool)
	defer close(out)
	out <- false
	return out
}
