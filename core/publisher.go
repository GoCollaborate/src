package core

import (
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server"
	"github.com/GoCollaborate/server/task"
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
	LocalTasks   []func() task.Task
	ExposedTasks []func() task.Task
}

type Workable interface {
	Attach() uint64
	BatchAttach(amount int) []uint64
	Detach(w *server.Worker)
	LaunchAll() error
	Launch(ID uint64) error
	Enqueue(t task.Task)
	EnqueueMulti(ts []task.Task)
	CountTasks() []int
	CountWorkers() int
	Close() bool
}

var singleton *Publisher
var once sync.Once

func GetPublisherInstance(lg *logger.Logger) *Publisher {
	once.Do(func() {
		singleton = &Publisher{new(server.Master), lg, *new([]func() task.Task), *new([]func() task.Task)}
		// initialise your custom tasks below
		singleton.AddLocal(TaskA,
			TaskB,
			TaskC).AddExposed(TaskA,
			TaskB,
			TaskC)
	})
	return singleton
}

func (p *Publisher) Connect(w Workable) {
	p.Workable = w
}

func (p *Publisher) Distribute(tsk task.Task) {
	p.Workable.Enqueue(tsk)
}

func (p *Publisher) DistributeGroup(tsks ...task.Task) {
	p.Workable.EnqueueMulti(tsks)
}

func (p *Publisher) AddLocal(tsks ...func() task.Task) *Publisher {
	p.LocalTasks = append(p.LocalTasks, tsks...)
	return p
}

func (p *Publisher) AddExposed(tsks ...func() task.Task) *Publisher {
	p.ExposedTasks = append(p.ExposedTasks, tsks...)
	return p
}

func Delay(sec time.Duration) {
	tm := time.NewTimer(sec * time.Second)
	<-tm.C
}
