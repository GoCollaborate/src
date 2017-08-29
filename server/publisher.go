package server

import (
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server/servershared"
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
	Workable    Workable
	Logger      *logger.Logger
	LocalTasks  []func() (interface{}, []string)
	SharedTasks []func() (interface{}, []string)
}

type Workable interface {
	Attach() uint64
	BatchAttach(amount int) []uint64
	Detach(w *servershared.Worker)
	LaunchAll() error
	Launch(ID uint64) error
	Enqueue(t ...task.Task)
	Proceed(t *task.Task) error
	Done(...task.Task) error
	CountTasks() []int
	CountWorkers() int
	Close() bool
}

var singleton *Publisher
var once sync.Once

func GetPublisherInstance(lg *logger.Logger) *Publisher {
	once.Do(func() {
		singleton = &Publisher{Dummy(), lg, *new([]func() (interface{}, []string)), *new([]func() (interface{}, []string))}
	})
	return singleton
}

func (p *Publisher) Connect(w Workable) {
	p.Workable = w
}

func (p *Publisher) LocalDistribute(tsks ...task.Task) {
	p.Workable.Enqueue(tsks...)
}

func (p *Publisher) SharedDistribute(tsks ...*task.Task) {
	for _, t := range tsks {
		p.Workable.Proceed(t)
	}
}

func (p *Publisher) SyncDistribute(tsks ...task.Task) error {
	return p.Workable.Done(tsks...)
}

func (p *Publisher) AddLocal(tsks ...func() (interface{}, []string)) *Publisher {
	p.LocalTasks = append(p.LocalTasks, tsks...)
	return p
}

func (p *Publisher) AddShared(tsks ...func() (interface{}, []string)) *Publisher {
	p.SharedTasks = append(p.SharedTasks, tsks...)
	return p
}

func (p *Publisher) HandleLocal(router *mux.Router, api string, tskFunc func() (interface{}, []string)) *Publisher {
	_tskFunc, methods := tskFunc()
	fun := _tskFunc.(func(w http.ResponseWriter, r *http.Request) task.Task)
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		p.LocalDistribute(fun(w, r))
	}).Methods(methods...)
	return p
}

func (p *Publisher) HandleShared(router *mux.Router, api string, tskFunc func() (interface{}, []string)) *Publisher {
	_tskFunc, methods := tskFunc()
	fun := _tskFunc.(func(w http.ResponseWriter, r *http.Request) task.Task)
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		f := fun(w, r)
		p.SharedDistribute(&f)
	}).Methods(methods...)
	return p
}

func Delay(sec time.Duration) {
	tm := time.NewTimer(sec * time.Second)
	<-tm.C
}

type dummyWorkable struct {
}

func Dummy() *dummyWorkable {
	return new(dummyWorkable)
}

func (d *dummyWorkable) Attach() uint64 {
	return 0
}

func (d *dummyWorkable) BatchAttach(amount int) []uint64 {
	return []uint64{}
}

func (d *dummyWorkable) Detach(w *servershared.Worker) {
	return
}

func (d *dummyWorkable) LaunchAll() error {
	return nil
}

func (d *dummyWorkable) Launch(ID uint64) error {
	return nil
}

func (d *dummyWorkable) Enqueue(t ...task.Task) {
	return
}

func (d *dummyWorkable) Done(...task.Task) error {
	return nil
}

func (d *dummyWorkable) CountTasks() []int {
	return []int{}
}

func (d *dummyWorkable) CountWorkers() int {
	return 0
}

func (d *dummyWorkable) Close() bool {
	return false
}

func (d *dummyWorkable) Proceed(tsk *task.Task) error {
	return nil
}
