package workable

import (
	"errors"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote/collaborator"
	"github.com/GoCollaborate/server/mapper"
	"github.com/GoCollaborate/server/reducer"
	"github.com/GoCollaborate/server/servershared"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/store"
	"strconv"
	"sync"
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
	bookkeeper  *collaborator.BookKeeper
}

func NewMaster(bkp *collaborator.BookKeeper, args ...*logger.Logger) *Master {
	if len(args) > 0 {
		return &Master{0, make(map[uint64]*servershared.Worker), args[0], make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), bkp}
	}
	return &Master{0, make(map[uint64]*servershared.Worker), nil, make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), make(chan *task.Task), bkp}
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

// proceed will process the parallel tasks at the same stage
func (m *Master) Proceed(tsks ...*task.Task) (*sync.WaitGroup, chan error) {

	var (
		NumTsks = len(tsks)
		errchan = make(chan error, NumTsks)
		wg      = &sync.WaitGroup{}
	)

	wg.Add(NumTsks)

	for _, tsk := range tsks {
		go func() {
			defer wg.Done()
			var (
				mp   mapper.Mapper
				rd   reducer.Reducer
				maps map[int64]*task.Task
				err  error
			)
			fs := store.GetInstance()
			mp, err = fs.GetMapper(tsk.Mapper)

			if err != nil {
				errchan <- err
				return
			}

			maps, err = mp.Map(tsk)

			if err != nil {
				errchan <- err
				return
			}

			maps, err = m.bookkeeper.SyncDistribute(maps)

			if err != nil {
				errchan <- err
				return
			}

			rd, err = fs.GetReducer(tsk.Reducer)

			if err != nil {
				errchan <- err
				return
			}

			errchan <- rd.Reduce(maps, tsk)
			return
		}()
	}

	return wg, errchan
}

// sequentially execute all tasks
func (m *Master) Done(ts ...*task.Task) error {
	fs := store.GetInstance()
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
