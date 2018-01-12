package worker

import (
	"fmt"
	"github.com/GoCollaborate/src/artifacts/task"
	"github.com/GoCollaborate/src/logger"
	"github.com/GoCollaborate/src/store"
)

type Worker struct {
	ID          uint
	Alive       bool
	BaseTasks   chan *task.TaskFuture
	LowTasks    chan *task.TaskFuture
	MediumTasks chan *task.TaskFuture
	HighTasks   chan *task.TaskFuture
	UrgentTasks chan *task.TaskFuture
	Exit        chan bool
}

func (w *Worker) Start() {
	fs := store.GetInstance()
	go func() {
		for {
			select {
			case <-w.Exit:
				return
			default:
				tkf := preselect(
					w.UrgentTasks,
					w.HighTasks,
					w.MediumTasks,
					w.LowTasks,
					w.BaseTasks,
				)
				tk := tkf.Receive()
				logger.LogNormal(fmt.Sprintf(
					"Worker%v:, Task Level:%v",
					w.ID,
					tk.Priority,
				))
				logger.GetLoggerInstance().
					LogNormal(fmt.Sprintf(
						"Worker%v:, Task Level:%v",
						w.ID,
						tk.Priority,
					))
				tkf.Return(fs.Call(
					(*tk).Consumable,
					&(*tk).Source,
					&(*tk).Result,
					(*tk).Context,
				))
			}
		}
	}()
}

func (w *Worker) GetID() uint {
	return w.ID
}

func (w *Worker) Quit() {
	w.Exit <- true
}

func preselect(a, b, c, d, e chan *task.TaskFuture) *task.TaskFuture {
	select {
	case x := <-a:
		return x
	default:
	}

	select {
	case x := <-a:
		return x
	case x := <-b:
		return x
	default:
	}

	select {
	case x := <-a:
		return x
	case x := <-b:
		return x
	case x := <-c:
		return x
	default:
	}

	select {
	case x := <-a:
		return x
	case x := <-b:
		return x
	case x := <-c:
		return x
	case x := <-d:
		return x
	default:
	}

	select {
	case x := <-a:
		return x
	case x := <-b:
		return x
	case x := <-c:
		return x
	case x := <-d:
		return x
	case x := <-e:
		return x
	}
}
