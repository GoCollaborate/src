package master

import (
	"errors"
	"github.com/GoCollaborate/src/artifacts/task"
	"github.com/GoCollaborate/src/artifacts/worker"
	"github.com/GoCollaborate/src/constants"
	"strconv"
	"time"
)

type Master struct {
	Count       uint
	List        map[uint]*worker.Worker
	BaseTasks   chan *task.TaskFuture
	LowTasks    chan *task.TaskFuture
	MediumTasks chan *task.TaskFuture
	HighTasks   chan *task.TaskFuture
	UrgentTasks chan *task.TaskFuture
}

func NewMaster() *Master {
	return &Master{0, make(map[uint]*worker.Worker), make(chan *task.TaskFuture), make(chan *task.TaskFuture), make(chan *task.TaskFuture), make(chan *task.TaskFuture), make(chan *task.TaskFuture)}
}

func (m *Master) Enqueue(ts map[int]*task.Task) map[int]*task.TaskFuture {
	futures := make(map[int]*task.TaskFuture)
	for k, t := range ts {
		futures[k] = task.NewTaskFuture(t)
		switch t.Priority.GetPriority() {
		case task.URGENT:
			m.UrgentTasks <- futures[k]
		case task.HIGH:
			m.HighTasks <- futures[k]
		case task.MEDIUM:
			m.MediumTasks <- futures[k]
		case task.LOW:
			m.LowTasks <- futures[k]
		default:
			m.BaseTasks <- futures[k]
		}
	}
	return futures
}

func (m *Master) DoneMulti(ts map[int]*task.Task) error {
	futures := m.Enqueue(ts)
	for _, f := range futures {
		select {
		case <-f.IsDone():
			f.Close()
			continue
		case <-time.After(constants.DefaultTaskExpireTime):
			return constants.ErrTimeout
		}
	}
	return nil
}

func (m *Master) Done(t *task.Task) error {
	future := task.NewTaskFuture(t)
	switch t.Priority.GetPriority() {
	case task.URGENT:
		m.UrgentTasks <- future
	case task.HIGH:
		m.HighTasks <- future
	case task.MEDIUM:
		m.MediumTasks <- future
	case task.LOW:
		m.LowTasks <- future
	default:
		m.BaseTasks <- future
	}

	select {
	case <-future.IsDone():
		future.Close()
		return nil
	case <-time.After(constants.DefaultTaskExpireTime):
		return constants.ErrTimeout
	}

	return nil
}

func (m *Master) Attach() uint {
	w := &worker.Worker{m.Count, true, m.BaseTasks, m.LowTasks, m.MediumTasks, m.HighTasks, m.UrgentTasks, make(chan bool)}
	m.Count++
	m.List[w.ID] = w
	return w.ID
}

func (m *Master) BatchAttach(amount int) []uint {
	var ids []uint = []uint{}
	for i := 0; i < amount; i++ {
		w := &worker.Worker{m.Count, true, m.BaseTasks, m.LowTasks, m.MediumTasks, m.HighTasks, m.UrgentTasks, make(chan bool)}
		m.Count++
		m.List[w.ID] = w
		ids = append(ids, w.ID)
	}
	return ids
}

func (m *Master) Detach(w *worker.Worker) {
	delete(m.List, w.ID)
}

func (m *Master) LaunchAll() error {
	for i, wk := range m.List {
		if wk.Alive {
			wk.Start()
			continue
		}
		return errors.New("Worker ID Not Exist Error:" + strconv.Itoa(int(i)))
	}
	return nil
}

func (m *Master) Launch(ID uint) error {
	if wk := m.List[ID]; wk.Alive {
		wk.Start()
		return nil
	}
	return errors.New("Worker ID Not Exist Error:" + strconv.Itoa(int(ID)))
}

// count the number of tasks that the are queuing in channel
func (m *Master) CountTasks() []int {
	return []int{len(m.BaseTasks), len(m.LowTasks), len(m.MediumTasks), len(m.HighTasks), len(m.UrgentTasks)}
}

// count the number of workers that are governed by the current master
func (m *Master) CountWorkers() int {
	return len(m.List)
}

func (m *Master) Close() bool {
	close(m.BaseTasks)
	close(m.LowTasks)
	close(m.MediumTasks)
	close(m.HighTasks)
	close(m.UrgentTasks)
	return true
}
