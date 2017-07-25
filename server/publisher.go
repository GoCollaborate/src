package server

import (
	"github.com/GoCollaborate/logger"
	"time"
)

const (
	MAXIDLECONNECTIONS int           = 20
	REQUESTTIMEOUT     time.Duration = 5
	UPDATEINTERVAL     time.Duration = 1
)

type Publisher struct {
	Workable Workable
	Logger   *logger.Logger
}

type Workable interface {
	Attach() uint64
	BatchAttach(amount int) []uint64
	Detach(w *Worker)
	LaunchAll() error
	Launch(ID uint64) error
	Enqueue(t Task)
	EnqueueMulti(ts []Task)
	CountTasks() []int
	CountWorkers() int
	Close() bool
}

func NewPublisher(lg *logger.Logger) *Publisher {
	return &Publisher{new(Master), lg}
}

func (p *Publisher) Connect(w Workable) {
	p.Workable = w
}

func (p *Publisher) Distribute() {
	p.Workable.Enqueue(TaskA())
	p.Workable.Enqueue(TaskB())
	p.Workable.Enqueue(TaskC())
}

func (p *Publisher) DistributeGroup() {
	for {
		p.Workable.EnqueueMulti([]Task{TaskA(), TaskB(), TaskC(), TaskA(), TaskB(), TaskC(), TaskA(), TaskB(), TaskC(), TaskA(), TaskB(), TaskC(), TaskA(), TaskB(), TaskC(), TaskA(), TaskB(), TaskC(), TaskA(), TaskB(), TaskC(), TaskA(), TaskB(), TaskC()})
		Delay(3)
	}
}

func Delay(sec time.Duration) {
	tm := time.NewTimer(sec * time.Second)
	<-tm.C
}
