package iworkable

import (
	"github.com/GoCollaborate/artifacts/task"
	"github.com/GoCollaborate/artifacts/worker"
	"sync"
)

type Workable interface {
	Attach() uint
	BatchAttach(amount int) []uint
	Detach(w *worker.Worker)
	LaunchAll() error
	Launch(ID uint) error
	Enqueue(t map[int]*task.Task)
	Done(*task.Task) error
	DoneMulti(tsks map[int]*task.Task) error
	CountTasks() []int
	CountWorkers() int
	Close() bool
}

type dummyWorkable struct {
}

func Dummy() *dummyWorkable {
	return new(dummyWorkable)
}

func (d *dummyWorkable) Attach() uint {
	return 0
}

func (d *dummyWorkable) BatchAttach(amount int) []uint {
	return []uint{}
}

func (d *dummyWorkable) Detach(w *worker.Worker) {
	return
}

func (d *dummyWorkable) LaunchAll() error {
	return nil
}

func (d *dummyWorkable) Launch(ID uint) error {
	return nil
}

func (d *dummyWorkable) Enqueue(t map[int]*task.Task) {
	return
}

func (d *dummyWorkable) Done(*task.Task) error {
	return nil
}

func (d *dummyWorkable) DoneMulti(tsks map[int]*task.Task) error {
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

func (d *dummyWorkable) Proceed(tsks map[int]*task.Task) (*sync.WaitGroup, chan error) {
	return &sync.WaitGroup{}, nil
}
