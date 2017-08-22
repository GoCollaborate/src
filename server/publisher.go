package server

import (
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server/task"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)

const (
	MAXIDLECONNECTIONS int           = 20
	REQUESTTIMEOUT     time.Duration = 5
	UPDATEINTERVAL     time.Duration = 1
)

type Publisher struct {
	Workable     Workable
	Logger       *logger.Logger
	LocalTasks   []func() (interface{}, []string)
	ExposedTasks []func() (interface{}, []string)
}

type Workable interface {
	Attach() uint64
	BatchAttach(amount int) []uint64
	Detach(w *Worker)
	LaunchAll() error
	Launch(ID uint64) error
	Enqueue(t ...task.Task)
	CountTasks() []int
	CountWorkers() int
	Close() bool
}

var singleton *Publisher
var once sync.Once

func GetPublisherInstance(lg *logger.Logger) *Publisher {
	once.Do(func() {
		singleton = &Publisher{new(Master), lg, *new([]func() (interface{}, []string)), *new([]func() (interface{}, []string))}
	})
	return singleton
}

func (p *Publisher) Connect(w Workable) {
	p.Workable = w
}

func (p *Publisher) Distribute(tsks ...task.Task) {
	p.Workable.Enqueue(tsks...)
}

func (p *Publisher) AddLocal(tsks ...func() (interface{}, []string)) *Publisher {
	p.LocalTasks = append(p.LocalTasks, tsks...)
	return p
}

func (p *Publisher) AddExposed(tsks ...func() (interface{}, []string)) *Publisher {
	p.ExposedTasks = append(p.ExposedTasks, tsks...)
	return p
}

func (p *Publisher) Handle(router *mux.Router, api string, tskFunc func() (interface{}, []string)) *Publisher {
	_tskFunc, methods := tskFunc()
	fun := _tskFunc.(func(w http.ResponseWriter, r *http.Request) task.Task)
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		p.Distribute(fun(w, r))
	}).Methods(methods...)
	return p
}

func Delay(sec time.Duration) {
	tm := time.NewTimer(sec * time.Second)
	<-tm.C
}
