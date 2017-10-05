package master

import (
	"errors"
	"github.com/GoCollaborate/artifacts/task"
	"github.com/GoCollaborate/artifacts/worker"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/store"
	"strconv"
	"time"
)

type Master struct {
	Count       uint
	List        map[uint]*worker.Worker
	BaseTasks   chan *task.Task
	LowTasks    chan *task.Task
	MediumTasks chan *task.Task
	HighTasks   chan *task.Task
	UrgentTasks chan *task.Task
}

func NewMaster() *Master {
	return &Master{0, make(map[uint]*worker.Worker), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task)}
}

func (m *Master) Enqueue(ts map[int]*task.Task) {
	for _, t := range ts {
		switch t.Priority.GetPriority() {
		case task.URGENT:
			m.UrgentTasks <- t
		case task.HIGH:
			m.HighTasks <- t
		case task.MEDIUM:
			m.MediumTasks <- t
		case task.LOW:
			m.LowTasks <- t
		default:
			m.BaseTasks <- t
		}
	}
}

// sequentially execute all tasks
func (m *Master) DoneMulti(tsks map[int]*task.Task) error {
	fs := store.GetInstance()
	for _, t := range tsks {
		switch t.Priority.GetPriority() {
		case task.URGENT:
			m.UrgentTasks <- t
		case task.HIGH:
			m.HighTasks <- t
		case task.MEDIUM:
			m.MediumTasks <- t
		case task.LOW:
			m.LowTasks <- t
		default:
			m.BaseTasks <- t
		}

		select {
		case <-fs.Listen(t.Consumable):
			continue
		case <-time.After(constants.DefaultTaskExpireTime):
			return constants.ErrTimeout
		}
	}
	return nil
}

// sequentially execute one task
func (m *Master) Done(t *task.Task) error {
	fs := store.GetInstance()
	switch t.Priority.GetPriority() {
	case task.URGENT:
		m.UrgentTasks <- t
	case task.HIGH:
		m.HighTasks <- t
	case task.MEDIUM:
		m.MediumTasks <- t
	case task.LOW:
		m.LowTasks <- t
	default:
		m.BaseTasks <- t
	}

	select {
	case <-fs.Listen(t.Consumable):
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
