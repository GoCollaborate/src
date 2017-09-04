package workable

import (
	"errors"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/funcstore"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote/collaborator"
	"github.com/GoCollaborate/server/mapper"
	"github.com/GoCollaborate/server/reducer"
	"github.com/GoCollaborate/server/servershared"
	"github.com/GoCollaborate/server/task"
	"strconv"
	"time"
)

type Master struct {
	Count       uint64
	List        map[uint64]*servershared.Worker
	Logger      *logger.Logger
	BaseTasks   chan *task.Task
	LowTasks    chan *task.Task
	MediumTasks chan *task.Task
	HighTasks   chan *task.Task
	UrgentTasks chan *task.Task
	mapper      mapper.Mapper
	reducer     reducer.Reducer
	bookkeeper  *collaborator.BookKeeper
}

func NewMaster(bkp *collaborator.BookKeeper, args ...*logger.Logger) *Master {
	if len(args) > 0 {
		return &Master{0, make(map[uint64]*servershared.Worker), args[0], make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), mapper.Default(), reducer.Default(), bkp}
	}
	return &Master{0, make(map[uint64]*servershared.Worker), nil, make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), mapper.Default(), reducer.Default(), bkp}
}

func (m *Master) Enqueue(ts ...*task.Task) {
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

func (m *Master) Mapper(mp mapper.Mapper) *Master {
	m.mapper = mp
	return m
}

func (m *Master) Reducer(rd reducer.Reducer) *Master {
	m.reducer = rd
	return m
}

func (m *Master) Proceed(tsk *task.Task) error {
	maps, err := m.mapper.Map(tsk)

	if err != nil {
		return err
	}

	maps, err = m.bookkeeper.SyncDistribute(maps)

	return m.reducer.Reduce(maps, tsk)
}

// sequentially execute all tasks
func (m *Master) Done(ts ...*task.Task) error {
	fs := funcstore.GetFSInstance()
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

		select {
		case <-fs.Listen(t.Consumable):
			continue
		case <-time.After(constants.DefaultTaskExpireTime):
			return constants.ErrTimeout
		}

	}
	return nil
}

func (m *Master) Attach() uint64 {
	w := &servershared.Worker{m.Count, true, m.BaseTasks, m.LowTasks, m.MediumTasks, m.HighTasks, m.UrgentTasks, make(chan bool)}
	m.Count++
	m.List[w.ID] = w
	return w.ID
}

func (m *Master) BatchAttach(amount int) []uint64 {
	var ids []uint64 = []uint64{}
	for i := 0; i < amount; i++ {
		w := &servershared.Worker{m.Count, true, m.BaseTasks, m.LowTasks, m.MediumTasks, m.HighTasks, m.UrgentTasks, make(chan bool)}
		m.Count++
		m.List[w.ID] = w
		ids = append(ids, w.ID)
	}
	return ids
}

func (m *Master) Detach(w *servershared.Worker) {
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
