package store

import (
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server/executor"
	"github.com/GoCollaborate/server/mapper"
	"github.com/GoCollaborate/server/reducer"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/utils"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)

var router *mux.Router
var singleton *FS
var once sync.Once
var onceRouter sync.Once
var mu sync.Mutex

type color int

const (
	white color = iota
	grey
	black
)

func GetRouter() *mux.Router {
	onceRouter.Do(func() {
		router = mux.NewRouter()
	})
	return router
}

func GetInstance() *FS {
	once.Do(func() {
		singleton = &FS{make(map[string]func(source *[]task.Countable,
			result *[]task.Countable,
			context *task.TaskContext) chan bool),
			make(map[string]chan bool),
			make(map[string]*color),
			make(map[string]executor.Executor),
			make(map[string]*task.Job),
			make(map[string]*JobFunc),
			make(map[string]*JobFunc)}
		singleton.sweep()
	})
	return singleton
}

type FS struct {
	Funcs map[string]func(source *[]task.Countable,
		result *[]task.Countable,
		context *task.TaskContext) chan bool
	Outbound   map[string]chan bool
	memstack   map[string]*color
	executors  map[string]executor.Executor
	jobs       map[string]*task.Job
	SharedJobs map[string]*JobFunc
	LocalJobs  map[string]*JobFunc
}

type JobFunc struct {
	F         func(w http.ResponseWriter, r *http.Request) *task.Job
	Methods   []string
	Signature string
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

func (fs *FS) SetMapper(mp mapper.Mapper, name string) {
	exe := executor.Default()
	exe.Todo(mp.Map)
	exe.Type(constants.ExecutorTypeMapper)
	fs.executors[name] = exe
}

func (fs *FS) SetReducer(rd reducer.Reducer, name string) {
	exe := executor.Default()
	exe.Todo(rd.Reduce)
	exe.Type(constants.ExecutorTypeReducer)
	fs.executors[name] = exe
}

func (fs *FS) GetExecutor(name string) (executor.Executor, error) {
	if exe := fs.executors[name]; exe != nil {
		return exe, nil
	}
	return executor.Default(), constants.ErrExecutorNotFound
}

func (fs *FS) SetJob(j *task.Job) {
	fs.jobs[j.Id()] = j
}

func (fs *FS) GetJob(id string) (*task.Job, error) {
	if j := fs.jobs[id]; j != nil {
		return j, nil
	}
	return task.MakeJob(), constants.ErrJobNotExist
}

func (fs *FS) SetShared(key string, val *JobFunc) {
	fs.SharedJobs[key] = val
}

func (fs *FS) SetLocal(key string, val *JobFunc) {
	fs.LocalJobs[key] = val
}

func (fs *FS) GetLocal(key string) (*JobFunc, error) {
	if j := fs.LocalJobs[key]; j != nil {
		return j, nil
	}
	return new(JobFunc), constants.ErrJobNotExist
}

func (fs *FS) GetShared(key string) (*JobFunc, error) {
	if j := fs.SharedJobs[key]; j != nil {
		return j, nil
	}
	return new(JobFunc), constants.ErrJobNotExist
}

func (fs *FS) AddLocal(methods []string, jobs ...func(w http.ResponseWriter, r *http.Request) *task.Job) {
	for _, f := range jobs {
		signature := utils.StripRouteToAPIRoute(utils.ReflectFuncName(f))
		fs.LocalJobs[signature] = &JobFunc{f, methods, signature}
	}
}

func (fs *FS) AddShared(methods []string, jobs ...func(w http.ResponseWriter, r *http.Request) *task.Job) {
	for _, f := range jobs {
		signature := utils.StripRouteToAPIRoute(utils.ReflectFuncName(f))
		fs.SharedJobs[signature] = &JobFunc{f, methods, signature}
	}
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
