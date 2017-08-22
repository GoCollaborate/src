package server

import (
	"errors"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server/task"
	"strconv"
)

type Master struct {
	Count       uint64
	List        map[uint64]*Worker
	Logger      *logger.Logger
	baseTasks   chan task.Task
	lowTasks    chan task.Task
	mediumTasks chan task.Task
	highTasks   chan task.Task
	urgentTasks chan task.Task
}

func NewMaster(args ...*logger.Logger) *Master {
	if len(args) > 0 {
		return &Master{0, make(map[uint64]*Worker), args[0], make(chan task.Task), make(chan task.Task), make(chan task.Task), make(chan task.Task), make(chan task.Task)}
	}
	return &Master{0, make(map[uint64]*Worker), nil, make(chan task.Task), make(chan task.Task), make(chan task.Task), make(chan task.Task), make(chan task.Task)}
}

func (m *Master) Enqueue(ts ...task.Task) {
	for _, t := range ts {
		switch t.Priority.GetPriority() {
		case task.URGENT:
			m.urgentTasks <- t
		case task.HIGH:
			m.highTasks <- t
		case task.MEDIUM:
			m.mediumTasks <- t
		case task.LOW:
			m.lowTasks <- t
		default:
			m.baseTasks <- t
		}
	}
}

func (m *Master) Attach() uint64 {
	w := &Worker{m.Count, *m, true, m.baseTasks, m.lowTasks, m.mediumTasks, m.highTasks, m.urgentTasks, make(chan bool)}
	m.Count++
	m.List[w.ID] = w
	return w.ID
}

func (m *Master) BatchAttach(amount int) []uint64 {
	var ids []uint64 = []uint64{}
	for i := 0; i < amount; i++ {
		w := &Worker{m.Count, *m, true, m.baseTasks, m.lowTasks, m.mediumTasks, m.highTasks, m.urgentTasks, make(chan bool)}
		m.Count++
		m.List[w.ID] = w
		ids = append(ids, w.ID)
	}
	return ids
}

func (m *Master) Detach(w *Worker) {
	delete(m.List, w.ID)
}

func (m *Master) LaunchAll() error {
	for i, wk := range m.List {
		if wk.Alive {
			wk.Start()
			continue
		}
		return errors.New("Worker_ID_Not_Exist_Error:" + strconv.Itoa(int(i)))
	}
	return nil
}

func (m *Master) Launch(ID uint64) error {
	if wk := m.List[ID]; wk.Alive {
		wk.Start()
		return nil
	}
	return errors.New("Worker_ID_Not_Exist_Error:" + strconv.Itoa(int(ID)))
}

// count the number of tasks that the are queuing in channel
func (m *Master) CountTasks() []int {
	return []int{len(m.baseTasks), len(m.lowTasks), len(m.mediumTasks), len(m.highTasks), len(m.urgentTasks)}
}

func (m *Master) CountWorkers() int {
	return len(m.List)
}

func (m *Master) Close() bool {
	close(m.baseTasks)
	close(m.lowTasks)
	close(m.mediumTasks)
	close(m.highTasks)
	close(m.urgentTasks)
	return true
}
