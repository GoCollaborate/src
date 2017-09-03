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
	LocalTasks  []TskFunc
	SharedTasks []TskFunc
}

type TskFunc struct {
	F       func(w http.ResponseWriter, r *http.Request) task.Task
	Methods []string
}

type Workable interface {
	Attach() uint64
	BatchAttach(amount int) []uint64
	Detach(w *servershared.Worker)
	LaunchAll() error
	Launch(ID uint64) error
	Enqueue(t ...*task.Task)
	Proceed(t *task.Task) error
	Done(...*task.Task) error
	CountTasks() []int
	CountWorkers() int
	Close() bool
}

var singleton *Publisher
var once sync.Once

func GetPublisherInstance() *Publisher {
	once.Do(func() {
		singleton = &Publisher{Dummy(), nil, *new([]TskFunc), *new([]TskFunc)}
	})
	return singleton
}

func Logger(lg *logger.Logger) {
	s := *singleton
	s.Logger = lg
}

func (p *Publisher) Connect(w Workable) {
	p.Workable = w
}

func (p *Publisher) LocalDistribute(tsks ...*task.Task) {
	p.Workable.Enqueue(tsks...)
}

func (p *Publisher) SharedDistribute(tsks ...*task.Task) {
	for _, t := range tsks {
		p.Workable.Proceed(t)
	}
}

func (p *Publisher) SyncDistribute(tsks ...*task.Task) chan *task.Task {
	ch := make(chan *task.Task)

	go func() {
		defer close(ch)
		err := p.Workable.Done(tsks...)
		if err != nil {
			logger.LogError("Execution Error:" + err.Error())
			ch <- &task.Task{}
		}
		for _, t := range tsks {
			ch <- t
		}
	}()
	return ch
}

func (p *Publisher) AddLocal(methods []string, tsks ...func(w http.ResponseWriter, r *http.Request) task.Task) *Publisher {
	tskFuncs := make([]TskFunc, len(tsks))
	for i, f := range tsks {
		tskFuncs[i] = TskFunc{f, methods}
	}
	p.LocalTasks = append(p.LocalTasks, tskFuncs...)
	return p
}

func (p *Publisher) AddShared(methods []string, tsks ...func(w http.ResponseWriter, r *http.Request) task.Task) *Publisher {
	tskFuncs := make([]TskFunc, len(tsks))
	for i, f := range tsks {
		tskFuncs[i] = TskFunc{f, methods}
	}
	p.SharedTasks = append(p.SharedTasks, tskFuncs...)
	return p
}

func (p *Publisher) HandleLocal(router *mux.Router, api string, tskFunc TskFunc) *Publisher {
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		t := tskFunc.F(w, r)
		p.LocalDistribute(&t)
	}).Methods(tskFunc.Methods...)
	return p
}

func (p *Publisher) HandleShared(router *mux.Router, api string, tskFunc TskFunc) *Publisher {
	router.HandleFunc(api, func(w http.ResponseWriter, r *http.Request) {
		t := tskFunc.F(w, r)
		p.SharedDistribute(&t)
	}).Methods(tskFunc.Methods...)
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

func (d *dummyWorkable) Enqueue(t ...*task.Task) {
	return
}

func (d *dummyWorkable) Done(...*task.Task) error {
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
