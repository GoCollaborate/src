package funcstore

import (
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/utils"
	"sync"
	"time"
)

var singleton *FS
var once sync.Once
var mu sync.Mutex

type color int

const (
	white color = iota
	grey
	black
)

func GetFSInstance() *FS {
	once.Do(func() {
		singleton = &FS{make(map[string]func(source *[]task.Countable,
			result *[]task.Countable,
			context *task.TaskContext) chan bool), make(map[string]chan bool), make(map[string]*color)}
		singleton.sweep()
	})
	return singleton
}

type FS struct {
	Funcs map[string]func(source *[]task.Countable,
		result *[]task.Countable,
		context *task.TaskContext) chan bool
	Outbound map[string]chan bool
	memstack map[string]*color
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

	mu.Lock()
	defer mu.Unlock()
	fs.Funcs[i] = f
	fs.Outbound[i] = make(chan bool)
}

func (fs *FS) HAdd(f func(source *[]task.Countable,
	result *[]task.Countable,
	context *task.TaskContext) chan bool) (hash string) {
	hash = utils.RandStringBytesMaskImprSrc(constants.DefaultHashLength)

	mu.Lock()
	defer mu.Unlock()
	fs.Funcs[hash] = f
	fs.Outbound[hash] = make(chan bool)
	*fs.memstack[hash] = grey
	return
}

func (fs *FS) Call(id string, source *[]task.Countable,
	result *[]task.Countable,
	context *task.TaskContext) {

	if f := fs.Funcs[id]; f != nil {
		if c := fs.memstack[id]; c != nil {
			fs.Outbound[id] <- <-f(source, result, context)
			*fs.memstack[id] = white
		}
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

func (fs *FS) sweep() {
	go func() {
		for {
			<-time.After(constants.DefaultGCInterval)
			// copy lookup table
			stack := fs.memstack
			for k, s := range stack {
				if *s == white {
					fs.delete(k)
				}
			}
		}
	}()
}

func (fs *FS) delete(id string) {
	mu.Lock()
	defer mu.Unlock()
	delete(fs.Funcs, id)
	delete(fs.Outbound, id)
	delete(fs.memstack, id)
}
